package postgres

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	History struct {
		Query     string `json:"query"`
		Timestamp string `json:"timestamp"`
	}

	Connection struct {
		client        *sqlx.DB
		lastQueryTime time.Time
		queryTimeout  time.Duration
		context       context.Context
		history       []History
		closed        bool
		logging       bool
	}

	Row map[string]interface{}

	Result struct {
		Rows    []Row    `json:"rows"`
		Columns []string `json:"columns"`
		Stats   *Stats   `json:"stats,omitempty"`
	}

	Stats struct {
		RowsCount       int       `json:"rows_count"`
		ColumnsCount    int       `json:"columns_count"`
		RowsAffected    int64     `json:"rows_affected"`
		QueryStartTime  time.Time `json:"query_start_time"`
		QueryFinishTime time.Time `json:"query_finish_time"`
		QueryDuration   int64     `json:"query_duration_ms"`
	}
)
