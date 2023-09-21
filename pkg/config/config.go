package config

import (
	"os"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

func IsDebug() bool {
	return env.Get("DEBUG", "false") == "true"
}

func IsLocal() bool {
	return env.Get("APP_ENV", "local") == "local"
}

func IsProduction() bool {
	return env.Get("APP_ENV", "local") == "production"
}

func IsStaging() bool {
	return env.Get("APP_ENV", "local") == "staging"
}

func AppEnv() string {
	return env.Get("APP_ENV", "local")
}

func TimeBrl() time.Time {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	return time.Now().In(loc)
}

func JwtExpiresIn() time.Duration {
	return time.Duration(env.GetInt("JWT_EXPIRES_IN_SECONDS", "86400"))
}

func JwtSecretKey() []byte {
	return []byte(env.Get("JWT_SECRET_KEY"))
}

func Pid() int {
	return os.Getpid()
}

func Hostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
