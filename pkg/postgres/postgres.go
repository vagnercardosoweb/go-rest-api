package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/slack"

	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewClient(ctx context.Context, logger *logger.Logger, config *Config) *Client {
	sslMode := "disable"
	if config.EnabledSSL {
		sslMode = "require"
	}

	dbx, err := sqlx.ConnectContext(
		ctx,
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s TimeZone=%s application_name=%s sslmode=%s search_path=%s",
			config.Host, config.Port, config.Username, config.Password, config.Database, config.Timezone, config.AppName, sslMode, config.Schema,
		),
	)

	if err != nil {
		panic(fmt.Errorf("failed to connect to postgres: %v", err))
	}

	dbx.SetMaxOpenConns(config.MaxOpenConn)
	dbx.SetConnMaxIdleTime(config.MaxIdleTimeConn)
	dbx.SetConnMaxLifetime(config.MaxLifetimeConn)
	dbx.SetMaxIdleConns(config.MaxIdleConn)

	client := &Client{
		dbx:         dbx,
		ctx:         ctx,
		afterCommit: make([]func(c *Client) error, 0),
		logger:      logger,
		config:      config,
		tx:          nil,
	}

	client.runMigrations()

	return client
}

func FromEnv(ctx context.Context, logger *logger.Logger) *Client {
	return NewClient(
		ctx,
		logger,
		&Config{
			Port:            env.GetAsInt("DB_PORT", "5432"),
			Host:            env.GetAsString("DB_HOST", "localhost"),
			Database:        env.GetAsString("DB_NAME", "development"),
			Username:        env.GetAsString("DB_USERNAME", "postgres"),
			Password:        env.GetAsString("DB_PASSWORD", "postgres"),
			Timezone:        env.GetAsString("DB_TIMEZONE", "UTC"),
			Schema:          env.GetAsString("DB_SCHEMA", "public"),
			AppName:         env.GetAsString("DB_APP_NAME", "app"),
			EnabledSSL:      env.GetAsBool("DB_ENABLED_SSL", "false"),
			MigrationDir:    env.GetAsString("DB_MIGRATION_DIR", "migrations"),
			AutoMigrate:     env.GetAsBool("DB_AUTO_MIGRATE", "false"),
			QueryTimeout:    time.Millisecond * time.Duration(env.GetAsInt("DB_QUERY_TIMEOUT", "7000")),
			MaxIdleTimeConn: time.Millisecond * time.Duration(env.GetAsInt("DB_MAX_IDLE_TIME_CONN", "15000")),
			MaxLifetimeConn: time.Millisecond * time.Duration(env.GetAsInt("DB_MAX_LIFETIME_CONN", "60000")),
			MaxOpenConn:     env.GetAsInt("DB_MAX_OPEN_CONN", "35"),
			MaxIdleConn:     env.GetAsInt("DB_MAX_IDLE_CONN", "0"),
			Logging:         env.GetAsBool("DB_LOGGING", "false"),
		},
	)
}

func (c *Client) withQueryTimeoutCtx() (context.Context, context.CancelFunc) {
	queryTimeout := c.config.QueryTimeout
	if queryTimeout > 0 {
		return context.WithTimeout(c.ctx, queryTimeout)
	}
	return c.ctx, func() {}
}

func (c *Client) DB() *sql.DB {
	return c.dbx.DB
}

func (c *Client) DBX() *sqlx.DB {
	return c.dbx
}

func (c *Client) Exec(query string, bind ...any) (sql.Result, error) {
	ctx, cancel := c.withQueryTimeoutCtx()
	defer cancel()

	var err error
	log := &Log{Query: query, Bind: bind, StartedAt: time.Now()}

	defer func() {
		c.log(log)
	}()

	var result sql.Result

	if c.tx != nil {
		result, err = c.tx.ExecContext(ctx, query, bind...)
	} else {
		result, err = c.dbx.ExecContext(ctx, query, bind...)
	}

	log.FinishedAt = time.Now()

	if err != nil {
		log.ErrorMessage = err.Error()
	}

	return result, err
}

func (c *Client) Query(dest any, query string, bind ...any) error {
	ctx, cancel := c.withQueryTimeoutCtx()
	defer cancel()

	var err error
	log := &Log{Query: query, Bind: bind, StartedAt: time.Now()}

	defer func() {
		c.log(log)
	}()

	if c.tx != nil {
		err = c.tx.SelectContext(ctx, dest, query, bind...)
	} else {
		err = c.dbx.SelectContext(ctx, dest, query, bind...)
	}

	log.FinishedAt = time.Now()

	if err != nil {
		log.ErrorMessage = err.Error()
	}

	return err
}

func (c *Client) QueryRow(dest any, query string, bind ...any) error {
	ctx, cancel := c.withQueryTimeoutCtx()
	defer cancel()

	var err error
	log := &Log{Query: query, Bind: bind, StartedAt: time.Now()}

	defer func() {
		c.log(log)
	}()

	if c.tx != nil {
		err = c.tx.GetContext(ctx, dest, query, bind...)
	} else {
		err = c.dbx.GetContext(ctx, dest, query, bind...)
	}

	log.FinishedAt = time.Now()

	if err != nil {
		log.ErrorMessage = err.Error()
	}

	return err
}

func (c *Client) TruncateTable(table string) error {
	_, err := c.Exec(fmt.Sprintf(`TRUNCATE TABLE "%s" RESTART IDENTITY CASCADE`, table))
	return err
}

func (c *Client) AfterCommit(fn func(client *Client) error) {
	if c.tx == nil {
		return
	}

	size := len(c.afterCommit)

	c.afterCommit = append(c.afterCommit, func(client *Client) error {
		err := fn(client)

		if err != nil {
			_ = slack.NewAlert().
				AddField("AppName", c.config.AppName, false).
				AddField("RequestId", c.logger.GetId(), false).
				AddError(fmt.Sprintf("ExecuteAfterCommitError[%d]", size), err).
				WithColor(slack.ColorError).
				Send()
		}

		return err
	})
}

func (c *Client) WithTx(fn func(*Client) (any, error)) (any, error) {
	tx, err := c.dbx.BeginTxx(c.ctx, nil)
	if err != nil {
		return nil, err
	}

	// prevents database lock in case of panic error
	defer func() {
		_ = tx.Rollback()
	}()

	client := c.Copy()
	client.tx = tx

	result, err := fn(client)

	if err != nil {
		if txError := tx.Rollback(); txError != nil {
			return nil, errors.New(errors.Input{
				RequestId: c.logger.GetId(),
				Message:   "DB_ROLLBACK_TRANSACTION_ERROR",
				Metadata: errors.Metadata{
					"txError": txError.Error(),
					"fnError": err.Error(),
				}},
			)
		}

		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	for _, fn := range client.afterCommit {
		go func() {
			_ = fn(client)
		}()
	}

	return result, nil
}

func (c *Client) Logger() *logger.Logger {
	return c.logger
}

func (c *Client) WithLogger(logger *logger.Logger) *Client {
	client := c.Copy()
	client.logger = logger

	return client
}

func (c *Client) Copy() *Client {
	return &Client{
		dbx:         c.dbx,
		ctx:         c.ctx,
		afterCommit: make([]func(client *Client) error, 0),
		logger:      c.logger,
		config:      c.config,
	}
}

func (c *Client) Close() error {
	return c.dbx.Close()
}

func (c *Client) LastLog() *Log {
	return c.lastLog
}

func (c *Client) Ping() error {
	return c.dbx.PingContext(c.ctx)
}
