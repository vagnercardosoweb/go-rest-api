package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Client struct {
	tx            *sqlx.Tx
	db            *sqlx.DB
	lastQueryTime time.Time
	context       context.Context
	logger        *logger.Logger
	history       []History
	config        *config
}

func NewClient(ctx context.Context, logger *logger.Logger, envPrefix EnvPrefix) *Client {
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

	return &Client{
		tx:      nil,
		db:      db,
		context: ctx,
		history: make([]History, 0),
		logger:  logger,
		config:  config,
	}
}

func (c *Client) withQueryTimeoutCtx() (context.Context, context.CancelFunc) {
	queryTimeout := c.config.getQueryTimeout()
	if queryTimeout > 0 {
		return context.WithTimeout(c.context, queryTimeout)
	}
	return c.context, func() {}
}

func (c *Client) GetDB() *sql.DB {
	return c.db.DB
}

func (c *Client) Exec(query string, bind ...any) (sql.Result, error) {
	query = c.normalizeQuery(query)

	ctx, cancel := c.withQueryTimeoutCtx()
	defer cancel()

	var err error
	history := History{Query: query, Bind: bind}

	defer func() {
		c.lastQueryTime = time.Now().UTC()
		history.Duration = history.FinishedAt.Sub(history.StartedAt).String()
		c.addHistory(history)
	}()

	var result sql.Result
	history.StartedAt = time.Now()

	if c.tx != nil {
		result, err = c.tx.ExecContext(ctx, query, bind...)
	} else {
		result, err = c.db.ExecContext(ctx, query, bind...)
	}

	history.FinishedAt = time.Now()

	if err != nil {
		history.ErrorMessage = err.Error()
	}

	return result, err
}

func (c *Client) Query(dest any, query string, bind ...any) error {
	query = c.normalizeQuery(query)

	ctx, cancel := c.withQueryTimeoutCtx()
	defer cancel()

	var err error
	history := History{Query: query, Bind: bind}

	defer func() {
		c.lastQueryTime = time.Now().UTC()
		history.Duration = history.FinishedAt.Sub(history.StartedAt).String()
		c.addHistory(history)
	}()

	history.StartedAt = time.Now()

	if c.tx != nil {
		err = c.tx.SelectContext(ctx, dest, query, bind...)
	} else {
		err = c.db.SelectContext(ctx, dest, query, bind...)
	}

	history.FinishedAt = time.Now()

	if err != nil {
		history.ErrorMessage = err.Error()
	}

	return err
}

func (c *Client) QueryOne(dest any, query string, bind ...any) error {
	query = c.normalizeQuery(query)

	ctx, cancel := c.withQueryTimeoutCtx()
	defer cancel()

	var err error
	history := History{Query: query, Bind: bind}

	defer func() {
		c.lastQueryTime = time.Now().UTC()
		history.Duration = history.FinishedAt.Sub(history.StartedAt).String()
		c.addHistory(history)
	}()

	history.StartedAt = time.Now()

	if c.tx != nil {
		err = c.tx.GetContext(ctx, dest, query, bind...)
	} else {
		err = c.db.GetContext(ctx, dest, query, bind...)
	}

	history.FinishedAt = time.Now()

	if err != nil {
		history.ErrorMessage = err.Error()
	}

	return err
}

func (c *Client) WithTx(fn func(*Client) (any, error)) (any, error) {
	tx, err := c.db.BeginTxx(c.context, nil)
	if err != nil {
		return nil, err
	}

	newConnection := c.Copy()
	newConnection.tx = tx

	result, err := fn(newConnection)

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

	return result, tx.Commit()
}

func (c *Client) WithLogger(logger *logger.Logger) *Client {
	newConnection := c.Copy()
	newConnection.logger = logger
	return newConnection
}

func (c *Client) GetLogger() *logger.Logger {
	return c.logger
}

func (c *Client) Copy() *Client {
	return &Client{
		db:      c.db,
		context: c.context,
		logger:  c.logger,
		history: make([]History, 0),
		config:  c.config,
	}
}

func (c *Client) normalizeQuery(query string) string {
	q := strings.TrimSpace(query)
	q = regexp.MustCompile(`\s+|\n|\t`).ReplaceAllString(q, " ")
	q = regexp.MustCompile(`\(\s`).ReplaceAllString(q, "(")
	q = regexp.MustCompile(`\s\)`).ReplaceAllString(q, ")")
	return q
}

func (c *Client) Close() error {
	return c.db.Close()
}

func (c *Client) LastQueryTime() time.Time {
	return c.lastQueryTime
}

func (c *Client) Ping() error {
	return c.db.PingContext(c.context)
}
