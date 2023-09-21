package postgres

import (
	"fmt"
	"strconv"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

type EnvPrefix string

type config struct {
	Port       int
	Host       string
	Database   string
	Username   string
	Password   string
	Timezone   string
	Schema     string
	AppName    string
	EnabledSSL bool
	Prefix     EnvPrefix
	Logging    bool
}

const (
	Default EnvPrefix = "DB"
)

func fromEnvPrefix(envPrefix EnvPrefix) *config {
	config := &config{}

	config.Prefix = envPrefix
	config.Port = config.getValueFromEnvToInt("PORT", 5432)
	config.Host = config.env("HOST")
	config.Database = config.env("NAME")
	config.Username = config.env("USERNAME")
	config.Password = config.env("PASSWORD")
	config.Timezone = config.env("TIMEZONE", "UTC")
	config.Schema = config.env("SCHEMA", "public")
	config.AppName = config.env("APP_NAME", "golang")
	config.EnabledSSL = config.env("ENABLED_SSL", "false") == "true"
	config.Logging = config.env("LOGGING", "false") == "true"

	return config
}

func (c *config) env(name string, defaultValue ...string) string {
	return env.Get(fmt.Sprintf("%s_%s", c.Prefix, name), defaultValue...)
}

func (c *config) getValueFromEnvToInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(c.env(key))
	if err != nil {
		return defaultValue
	}
	return value
}

func (c *config) getMaxPool() int {
	return c.getValueFromEnvToInt("POOL_MAX", 35)
}

func (c *config) getMaxIdleConn() int {
	return c.getValueFromEnvToInt("MAX_IDLE_CONN", 5)
}

func (c *config) getQueryTimeout() time.Duration {
	timeout := c.getValueFromEnvToInt("QUERY_TIMEOUT", 60)
	return time.Second * time.Duration(timeout)
}

func (c *config) getConnMaxLifetime() time.Duration {
	lifetime := c.getValueFromEnvToInt("MAX_LIFETIME_CONN", 60)
	return time.Second * time.Duration(lifetime)
}

func (c *config) getConnMaxIdleTime() time.Duration {
	idleTime := c.getValueFromEnvToInt("MAX_IDLE_TIME_CONN", 15)
	return time.Second * time.Duration(idleTime)
}
