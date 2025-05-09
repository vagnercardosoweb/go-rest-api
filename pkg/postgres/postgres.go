package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
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

	db, err := sqlx.ConnectContext(
		ctx,
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s TimeZone=%s application_name=%s sslmode=%s search_path=%s",
			config.Host, config.Port, config.Username, config.Password, config.Database, config.Timezone, config.AppName, sslMode, config.Schema,
		),
	)

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(config.MaxOpenConn)
	db.SetConnMaxIdleTime(config.MaxIdleTimeConn)
	db.SetConnMaxLifetime(config.MaxLifetimeConn)
	db.SetMaxIdleConns(config.MaxIdleConn)

	return &Client{
		db:          db,
		ctx:         ctx,
		afterCommit: make([]func(client *Client) error, 0),
		logger:      logger,
		config:      config,
		tx:          nil,
	}
}

func NewFromEnv(ctx context.Context, logger *logger.Logger) *Client {
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
			EnabledSSL:      env.GetAsString("DB_ENABLED_SSL", "false") == "true",
			QueryTimeout:    time.Millisecond * time.Duration(env.GetAsInt("DB_QUERY_TIMEOUT", "7000")),
			MaxIdleTimeConn: time.Millisecond * time.Duration(env.GetAsInt("DB_MAX_IDLE_TIME_CONN", "15000")),
			MaxLifetimeConn: time.Millisecond * time.Duration(env.GetAsInt("DB_MAX_LIFETIME_CONN", "60000")),
			MaxOpenConn:     env.GetAsInt("DB_MAX_OPEN_CONN", "35"),
			MaxIdleConn:     env.GetAsInt("DB_MAX_IDLE_CONN", "0"),
			Logging:         env.GetAsString("DB_LOGGING", "false") == "true",
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

func (c *Client) GetDb() *sql.DB {
	return c.db.DB
}

func (c *Client) Exec(query string, bind ...any) (sql.Result, error) {
	ctx, cancel := c.withQueryTimeoutCtx()
	defer cancel()

	var err error
	log := &Log{Query: query, Bind: bind}

	defer func() {
		c.log(log)
	}()

	log.StartedAt = time.Now()
	var result sql.Result

	if c.tx != nil {
		result, err = c.tx.ExecContext(ctx, query, bind...)
	} else {
		result, err = c.db.ExecContext(ctx, query, bind...)
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
	log := &Log{Query: query, Bind: bind}

	defer func() {
		c.log(log)
	}()

	log.StartedAt = time.Now()

	if c.tx != nil {
		err = c.tx.SelectContext(ctx, dest, query, bind...)
	} else {
		err = c.db.SelectContext(ctx, dest, query, bind...)
	}

	log.FinishedAt = time.Now()

	if err != nil {
		log.ErrorMessage = err.Error()
	}

	return err
}

func (c *Client) QueryOne(dest any, query string, bind ...any) error {
	ctx, cancel := c.withQueryTimeoutCtx()
	defer cancel()

	var err error
	log := &Log{Query: query, Bind: bind}

	defer func() {
		c.log(log)
	}()

	log.StartedAt = time.Now()

	if c.tx != nil {
		err = c.tx.GetContext(ctx, dest, query, bind...)
	} else {
		err = c.db.GetContext(ctx, dest, query, bind...)
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
				AddField("App Name", c.config.AppName, false).
				AddField("Request Id", c.logger.GetId(), false).
				AddError(fmt.Sprintf("ExecuteAfterCommitError[%d]", size), err).
				Send()
		}

		return err
	})
}

func (c *Client) WithTx(fn func(*Client) (any, error)) (any, error) {
	tx, err := c.db.BeginTxx(c.ctx, nil)
	if err != nil {
		return nil, err
	}

	// prevents database lock in case of panic error
	defer tx.Rollback()

	client := c.Copy()
	client.tx = tx

	result, err := fn(client)

	if err != nil {
		if txError := tx.Rollback(); txError != nil {
			return nil, errors.New(errors.Input{
				Message:    "Rollback Transaction Error",
				StatusCode: http.StatusInternalServerError,
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

func (c *Client) WithLogger(logger *logger.Logger) *Client {
	client := c.Copy()
	client.logger = logger
	return client
}

func (c *Client) GetLogger() *logger.Logger {
	return c.logger
}

func (c *Client) Copy() *Client {
	return &Client{
		db:          c.db,
		ctx:         c.ctx,
		afterCommit: make([]func(client *Client) error, 0),
		logger:      c.logger,
		config:      c.config,
	}
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) GetLastLog() *Log {
	return c.lastLog
}

func (c *Client) Ping() error {
	return c.db.PingContext(c.ctx)
}
