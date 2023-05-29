package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

type Connection struct {
	tx            *sqlx.Tx
	client        *sqlx.DB
	lastQueryTime time.Time
	context       context.Context
	history       []History
	config        *config
}

func Connect(ctx context.Context, envPrefix EnvPrefix) *Connection {
	config := fromEnvPrefix(envPrefix)

	sslMode := "disable"
	if config.EnabledSSL {
		sslMode = "require"
	}

	db, err := sqlx.ConnectContext(
		ctx,
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s TimeZone=%s application_name=%s sslmode=%s",
			config.Host, config.Port, config.Username, config.Password, config.Database, config.Timezone, config.AppName, sslMode,
		),
	)

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(config.getMaxPool())
	db.SetConnMaxIdleTime(config.getConnMaxIdleTime())
	db.SetConnMaxLifetime(config.getConnMaxLifetime())
	db.SetMaxIdleConns(config.getMaxIdleConn())

	return &Connection{
		tx:      nil,
		client:  db,
		context: ctx,
		history: make([]History, 0),
		config:  config,
	}
}

func (c *Connection) withQueryTimeoutCtx() (context.Context, context.CancelFunc) {
	queryTimeout := c.config.getQueryTimeout()
	if queryTimeout > 0 {
		return context.WithTimeout(c.context, queryTimeout)
	}
	return c.context, func() {}
}

func (c *Connection) GetDB() *sql.DB {
	return c.client.DB
}

func (c *Connection) Exec(query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := c.withQueryTimeoutCtx()
	defer cancel()

	var err error
	var result sql.Result
	history := History{
		Query:     query,
		CreatedAt: time.Now(),
		Arguments: args,
	}

	defer func() {
		c.lastQueryTime = time.Now().UTC()
		history.LatencyMs = history.FinishedAt.Sub(history.StartedAt).Milliseconds()
		c.addHistory(history)
	}()

	if c.config.Logging {
		//
	}

	history.StartedAt = time.Now()
	if c.tx != nil {
		result, err = c.tx.ExecContext(ctx, query, args...)
	} else {
		result, err = c.client.ExecContext(ctx, query, args...)
	}
	history.FinishedAt = time.Now()

	if err != nil {
		history.ErrorMessage = err.Error()
	}

	return result, err
}

func (c *Connection) Query(dest any, query string, args ...any) error {
	ctx, cancel := c.withQueryTimeoutCtx()
	defer cancel()

	var err error
	history := History{Query: query, Arguments: args, CreatedAt: time.Now()}

	defer func() {
		c.lastQueryTime = time.Now().UTC()
		history.LatencyMs = history.FinishedAt.Sub(history.StartedAt).Milliseconds()
		c.addHistory(history)
	}()

	if c.config.Logging {
		//
	}

	history.StartedAt = time.Now()
	if c.tx != nil {
		err = c.tx.SelectContext(ctx, dest, query, args...)
	} else {
		err = c.client.SelectContext(ctx, dest, query, args...)
	}
	history.FinishedAt = time.Now()

	if err != nil {
		history.ErrorMessage = err.Error()
	}

	return err
}

func (c *Connection) QueryOne(dest any, query string, args ...any) error {
	ctx, cancel := c.withQueryTimeoutCtx()
	defer cancel()

	var err error
	history := History{Query: query, Arguments: args, CreatedAt: time.Now()}

	defer func() {
		c.lastQueryTime = time.Now().UTC()
		history.LatencyMs = history.FinishedAt.Sub(history.StartedAt).Milliseconds()
		c.addHistory(history)
	}()

	if c.config.Logging {
		//
	}

	history.StartedAt = time.Now()
	if c.tx != nil {
		err = c.tx.GetContext(ctx, dest, query, args...)
	} else {
		err = c.client.GetContext(ctx, dest, query, args...)
	}
	history.FinishedAt = time.Now()

	if err != nil {
		history.ErrorMessage = err.Error()
	}

	return err
}

func (c *Connection) WithTx(fn func(*Connection) (any, error)) (any, error) {
	tx, err := c.client.BeginTxx(c.context, nil)
	if err != nil {
		return nil, err
	}

	result, err := fn(&Connection{
		tx:      tx,
		client:  c.client,
		context: c.context,
		history: make([]History, 0),
		config:  c.config,
	})

	if err != nil {
		if txError := tx.Rollback(); txError != nil {
			return nil, errors.New(errors.Input{
				Message:    "Rollback Transaction Error",
				StatusCode: http.StatusInternalServerError,
				Metadata: errors.Metadata{
					"fn_error": err.Error(),
					"tx_error": txError.Error(),
				}},
			)
		}

		return nil, err
	}

	return result, tx.Commit()
}

func (c *Connection) Close() error {
	return c.client.Close()
}

func (c *Connection) LastQueryTime() time.Time {
	return c.lastQueryTime
}

func (c *Connection) Ping() error {
	return c.client.PingContext(c.context)
}
