package postgres

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

type Config struct {
	Port         int
	Host         string
	Database     string
	Username     string
	Password     string
	Timezone     string
	Schema       string
	AppName      string
	MigrationDir string
	AutoMigrate  bool
	EnabledSSL   bool
	Logging      bool

	MaxIdleConn     int
	QueryTimeout    time.Duration
	MaxLifetimeConn time.Duration
	MaxIdleTimeConn time.Duration
	MaxOpenConn     int
}

type Client struct {
	db          *sqlx.DB
	tx          *sqlx.Tx
	config      *Config
	afterCommit []func(client *Client) error
	logger      *logger.Logger
	lastLog     *Log
	ctx         context.Context
}

type JsonToMap map[string]any

func (j *JsonToMap) Scan(value any) error {
	if value == nil {
		return nil
	}
	var data = value.([]byte)
	return json.Unmarshal(data, &j)
}

func (j JsonToMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type ArrayToMap []map[string]any

func (a *ArrayToMap) Scan(value any) error {
	if value == nil {
		return nil
	}
	var data = value.([]byte)
	return json.Unmarshal(data, &a)
}

func (a ArrayToMap) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}
	return json.Marshal(a)
}
