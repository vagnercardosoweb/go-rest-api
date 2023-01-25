package postgres

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"rest-api/shared"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewConnection(ctx context.Context) *Connection {
	port, _ := strconv.Atoi(shared.EnvRequiredByName("DB_PORT"))
	logging := shared.EnvGetByName("DB_LOGGING", "false") == "true"
	name := shared.EnvRequiredByName("DB_NAME")
	host := shared.EnvRequiredByName("DB_HOST")
	username := shared.EnvRequiredByName("DB_USERNAME")
	password := shared.EnvRequiredByName("DB_PASSWORD")
	timezone := shared.EnvRequiredByName("DB_TIMEZONE")

	db, err := sqlx.ConnectContext(
		ctx,
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
			host, port, username, password, name, timezone,
		),
	)
	if err != nil {
		log.Println("Database connection error")
		panic(err)
	}

	dbPoolMaxConnection, err := strconv.Atoi(shared.EnvGetByName("DB_POOL_MAX"))
	if err == nil && dbPoolMaxConnection > 0 {
		db.SetMaxOpenConns(dbPoolMaxConnection)
	}

	connection := &Connection{
		client:  db,
		context: ctx,
		history: make([]History, 0),
		logging: logging,
	}

	queryTimeout, _ := strconv.Atoi(shared.EnvGetByName("DB_QUERY_TIMEOUT", "0"))
	if queryTimeout > 0 {
		connection.queryTimeout = time.Second * time.Duration(queryTimeout)
	}

	return connection
}

func (connection *Connection) withQueryTimeoutCtx() (context.Context, context.CancelFunc) {
	if connection.queryTimeout > 0 {
		return context.WithTimeout(context.Background(), connection.queryTimeout)
	}
	return context.Background(), func() {}
}

func (connection *Connection) Exec(query string, args ...interface{}) (*Result, error) {
	ctx, cancel := connection.withQueryTimeoutCtx()
	defer cancel()

	queryStart := time.Now()
	res, err := connection.client.ExecContext(ctx, query, args...)
	queryFinish := time.Now()
	if err != nil {
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	result := Result{
		Rows:    []Row{{"result": affected}},
		Columns: []string{"result"},
		Stats: &Stats{
			RowsCount:       1,
			ColumnsCount:    1,
			QueryStartTime:  queryStart.UTC(),
			QueryFinishTime: queryFinish.UTC(),
			QueryDuration:   queryFinish.Sub(queryStart).Milliseconds(),
		},
	}

	return &result, nil
}

func (connection *Connection) Sqlx() *sqlx.DB {
	return connection.client
}

func (connection *Connection) Query(query string, args ...interface{}) (*Result, error) {
	if connection.client == nil {
		return nil, nil
	}

	defer func() {
		connection.lastQueryTime = time.Now().UTC()
		connection.addHistoryRecord(query)
	}()

	action := strings.ToLower(strings.Split(query, " ")[0])
	hasReturnValues := strings.Contains(strings.ToLower(query), " returning ")
	if (action == "update" || action == "delete") && !hasReturnValues {
		return connection.Exec(query, args...)
	}

	ctx, cancel := connection.withQueryTimeoutCtx()
	defer cancel()

	queryStart := time.Now()
	rows, err := connection.client.QueryxContext(ctx, query, args...)
	queryFinish := time.Now()
	if err != nil {
		log.Println("Failed query:", query, "\nArgs:", args)
		return nil, err
	}
	defer rows.Close()

	result := Result{
		Columns: []string{},
		Rows:    []Row{},
	}

	for rows.Next() {
		var row = Row{}
		rows.MapScan(row)
		for column, value := range row {
			result.Columns = append(result.Columns, column)
			if value == nil {
				row[column] = nil
			} else {
				if typeOf := reflect.TypeOf(value).Kind().String(); typeOf == "slice" {
					row[column] = string(value.([]byte))
				}
			}
		}
		result.Rows = append(result.Rows, row)
	}

	result.Stats = &Stats{
		RowsCount:       len(result.Rows),
		ColumnsCount:    len(result.Columns),
		QueryStartTime:  queryStart.UTC(),
		QueryFinishTime: queryFinish.UTC(),
		QueryDuration:   queryFinish.Sub(queryStart).Milliseconds(),
	}

	return &result, nil
}

// Close database connection
func (connection *Connection) Close() error {
	if connection.closed {
		return nil
	}

	defer func() {
		connection.closed = true
	}()

	if connection.client != nil {
		return connection.client.Close()
	}

	return nil
}

func (connection *Connection) IsClosed() bool {
	return connection.closed
}

func (connection *Connection) LastQueryTime() time.Time {
	return connection.lastQueryTime
}

func (connection *Connection) Test() error {
	return connection.client.Ping()
}

func (connection *Connection) addHistoryRecord(query string) {
	if !connection.hasHistoryRecord(query) {
		connection.history = append(connection.history, History{
			Query:     query,
			Timestamp: time.Now().String(),
		})
	}
}

func (connection *Connection) hasHistoryRecord(query string) bool {
	result := false

	for _, record := range connection.history {
		if record.Query == query {
			result = true
			break
		}
	}

	return result
}
