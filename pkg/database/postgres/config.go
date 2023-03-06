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

func getValueFromEnvToInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(env.Get(key))
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *Config) getMaxOpenConns() int {
	return getValueFromEnvToInt("DB_POOL_MAX", 50)
}

func (c *Config) getMaxIdleConns() int {
	return getValueFromEnvToInt("DB_MAX_IDLE_CONN", 30)
}

func (c *Config) getQueryTimeout() time.Duration {
	timeout := getValueFromEnvToInt("DB_QUERY_TIMEOUT", 3)
	return time.Second * time.Duration(timeout)
}

func (c *Config) getConnMaxLifetime() time.Duration {
	lifetime := getValueFromEnvToInt("DB_MAX_LIFETIME_CONN", 60)
	return time.Second * time.Duration(lifetime)
}

func (c *Config) getConnMaxIdleTime() time.Duration {
	idleTime := getValueFromEnvToInt("DB_MAX_IDLE_TIME_CONN", 15)
	return time.Second * time.Duration(idleTime)
}
