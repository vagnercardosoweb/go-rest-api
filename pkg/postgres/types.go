package postgres

import (
	"context"
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
