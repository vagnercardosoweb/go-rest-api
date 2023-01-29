package postgres

import (
	"strconv"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

type Config struct {
	Port     int
	Host     string
	Database string
	Username string
	Password string
	Timezone string
	Schema   string
	AppName  string
	Logging  bool
}

func newConfig() *Config {
	port, _ := strconv.Atoi(env.Get("DB_PORT", "5432"))
	logging := env.Get("DB_LOGGING", "false") == "true"

	return &Config{
		Port:     port,
		Host:     env.Required("DB_HOST"),
		Database: env.Required("DB_NAME"),
		Username: env.Required("DB_USERNAME"),
		Password: env.Required("DB_PASSWORD"),
		Timezone: env.Required("DB_TIMEZONE", "UTC"),
		Schema:   env.Required("SCHEMA", "public"),
		AppName:  env.Get("DB_APP_NAME", "go-structure"),
		Logging:  logging,
	}
}

func (c *Config) getMaxPoolConnection() int {
	max, err := strconv.Atoi(env.Get("DB_POOL_MAX", "5"))
	if err != nil || max <= 0 {
		return 5
	}
	return max
}

func (c *Config) getQueryTimeout() time.Duration {
	timeout, err := strconv.Atoi(env.Get("DB_QUERY_TIMEOUT", "0"))
	if err == nil && timeout > 0 {
		return time.Second * time.Duration(timeout)
	}
	return 0
}
