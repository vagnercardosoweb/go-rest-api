package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Connection struct {
	tx            *sqlx.Tx
	client        *sqlx.DB
	lastQueryTime time.Time
	context       context.Context
	history       []History
	config        *Config
}

func NewConnection(ctx context.Context) *Connection {
	config := newConfig()
	db, err := sqlx.ConnectContext(
		ctx,
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s application_name=%s",
			config.Host, config.Port, config.Username, config.Password, config.Database, config.Timezone, config.AppName,
		),
	)

	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(config.getMaxOpenConns())
	db.SetConnMaxIdleTime(config.getConnMaxIdleTime())
	db.SetConnMaxLifetime(config.getConnMaxLifetime())
	db.SetMaxIdleConns(config.getMaxIdleConns())

	return &Connection{
		tx:      nil,
		client:  db,
		context: ctx,
		history: make([]History, 0),
		config:  config,
	}
}

func (c *Connection) GetClient() *sqlx.DB {
	return c.client
}

func (c *Connection) withQueryTimeoutCtx() (context.Context, context.CancelFunc) {
	queryTimeout := c.config.getQueryTimeout()
	if queryTimeout > 0 {
		return context.WithTimeout(c.context, queryTimeout)
	}
	return c.context, func() {}
}

func (c *Connection) Exec(query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := c.withQueryTimeoutCtx()
	defer cancel()

	var err error
	var result sql.Result
	history := History{
		Query:     query,
		Arguments: args,
		CreatedAt: time.Now(),
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

func (c *Connection) Query(dest interface{}, query string, args ...interface{}) error {
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

func (c *Connection) WithTx(fn func(Connection) error) error {
	tx, err := c.client.BeginTxx(c.context, nil)
	if err != nil {
		return err
	}

	err = fn(Connection{
		tx:      tx,
		client:  c.client,
		context: c.context,
		history: make([]History, 0),
		config:  c.config,
	})

	if err != nil {
		if txError := tx.Rollback(); txError != nil {
			return errors.NewInternalServer(errors.Input{
				Message: "Rollback Transaction Error",
				Metadata: errors.Metadata{
					"fn_error": err.Error(),
					"tx_error": txError.Error(),
				}},
			)
		}

		return err
	}

	return tx.Commit()
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
