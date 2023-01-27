package postgres

import (
	"rest-api/shared"
	"strconv"
	"time"
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

func NewConfig() *Config {
	port, _ := strconv.Atoi(shared.EnvGetByName("DB_PORT", "5432"))
	logging := shared.EnvGetByName("DB_LOGGING", "false") == "true"

	return &Config{
		Port:     port,
		Host:     shared.EnvRequiredByName("DB_HOST"),
		Database: shared.EnvRequiredByName("DB_NAME"),
		Username: shared.EnvRequiredByName("DB_USERNAME"),
		Password: shared.EnvRequiredByName("DB_PASSWORD"),
		Timezone: shared.EnvRequiredByName("DB_TIMEZONE", "UTC"),
		Schema:   shared.EnvRequiredByName("SCHEMA", "public"),
		AppName:  shared.EnvGetByName("DB_APP_NAME", "go-structure"),
		Logging:  logging,
	}
}

func (c *Config) getMaxPoolConnection() int {
	max, err := strconv.Atoi(shared.EnvGetByName("DB_POOL_MAX", "5"))
	if err != nil || max <= 0 {
		return 5
	}
	return max
}

func (c *Config) getQueryTimeout() time.Duration {
	timeout, err := strconv.Atoi(shared.EnvGetByName("DB_QUERY_TIMEOUT", "0"))
	if err == nil && timeout > 0 {
		return time.Second * time.Duration(timeout)
	}
	return 0
}
