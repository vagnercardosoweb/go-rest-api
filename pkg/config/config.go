package config

import (
	"os"
	"strconv"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

var (
	AppEnv  = env.Get("APP_ENV", "local")
	IsDebug = env.Get("DEBUG", "false") == "true"

	Pid         = os.Getpid()
	Hostname, _ = os.Hostname()

	IsLocal      = AppEnv == "local"
	IsStaging    = AppEnv == "staging"
	IsProduction = AppEnv == "production"
)

const (
	PgConnectCtxKey    = "PgConnectCtxKey"
	RedisConnectCtxKey = "RedisConnectCtxKey"
	AuthTokenCtxKey    = "AuthTokenCtxKey"
	JwtPayloadCtxKey   = "JwtPayloadCtxKey"
	RequestIdCtxKey    = "RequestIdCtxKey"
)

func GetShutdownTimeout() time.Duration {
	timeout, err := strconv.Atoi(env.Get("SHUTDOWN_TIMEOUT", "0"))
	if err != nil {
		return 0
	}
	return time.Duration(timeout)
}
