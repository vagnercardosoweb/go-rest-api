package postgres

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"time"
)

type Config struct {
	Port       int
	Host       string
	Database   string
	Username   string
	Password   string
	Timezone   string
	Schema     string
	AppName    string
	EnabledSSL bool
	Logging    bool

	MaxIdleConn     int
	QueryTimeout    time.Duration
	MaxLifetimeConn time.Duration
	MaxIdleTimeConn time.Duration
	MaxOpenConn     int
}

type Client struct {
	db      *sqlx.DB
	tx      *sqlx.Tx
	lastLog *Log
	logger  *logger.Logger
	config  *Config
	ctx     context.Context
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
