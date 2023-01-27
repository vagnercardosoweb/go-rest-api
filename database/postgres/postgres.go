package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"rest-api/shared"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Connection struct {
	client        *sqlx.DB
	lastQueryTime time.Time
	context       context.Context
	history       []History
	config        *Config
	logger        *shared.Logger
}

func NewConnection(ctx context.Context) *Connection {
	config := NewConfig()
	db, err := sqlx.ConnectContext(
		ctx,
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s application_name=%s",
			config.Host, config.Port, config.Username, config.Password, config.Database, config.Timezone, config.AppName,
		),
	)

	if err != nil {
		log.Println("Database connection error")
		panic(err)
	}

	db.SetMaxOpenConns(config.getMaxPoolConnection())

	return &Connection{
		client:  db,
		context: ctx,
		history: make([]History, 0),
		logger:  shared.NewLogger(shared.Logger{Id: "POSTGRES"}),
		config:  config,
	}
}

func (connection *Connection) GetClient() *sqlx.DB {
	return connection.client
}

func (connection *Connection) withQueryTimeoutCtx() (context.Context, context.CancelFunc) {
	queryTimeout := connection.config.getQueryTimeout()
	if queryTimeout > 0 {
		return context.WithTimeout(connection.context, queryTimeout)
	}
	return connection.context, func() {}
}

func (connection *Connection) Exec(query string, args ...interface{}) (sql.Result, error) {
	ctx, cancel := connection.withQueryTimeoutCtx()
	defer cancel()

	result, err := connection.client.ExecContext(ctx, query, args...)

	if err != nil {
		connection.logger.
			AddMetadata("query", query).
			AddMetadata("arguments", args).
			Error(err.Error())

		return nil, err
	}

	return result, nil
}

func (connection *Connection) Query(dest interface{}, query string, args ...interface{}) error {
	ctx, cancel := connection.withQueryTimeoutCtx()
	defer cancel()

	history := History{Query: query, Arguments: args, CreatedAt: time.Now()}

	defer func() {
		connection.lastQueryTime = time.Now().UTC()

		history.LatencyMs = history.FinishedAt.Sub(history.StartedAt).Milliseconds()
		connection.addHistory(history)

		if history.ErrorMessage != "" {
			connection.logger.
				AddMetadata("query", query).
				AddMetadata("arguments", args).
				Error(history.ErrorMessage)
		}
	}()

	if connection.config.Logging {
		connection.logger.
			AddMetadata("query", query).
			AddMetadata("arguments", args).
			Debug("Executing query")
	}

	history.StartedAt = time.Now()
	err := connection.client.SelectContext(ctx, dest, query, args...)
	history.FinishedAt = time.Now()

	if err != nil {
		history.ErrorMessage = err.Error()
	}

	return err
}

func (connection *Connection) Close() error {
	connection.logger.Debug("Closing connection")
	return connection.client.Close()
}

func (connection *Connection) LastQueryTime() time.Time {
	return connection.lastQueryTime
}

func (connection *Connection) Ping() error {
	connection.logger.Debug("Ping connection")
	return connection.client.PingContext(connection.context)
}
