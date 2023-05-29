package config

import (
	"os"
	"strconv"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

var (
	Pid         = os.Getpid()
	Hostname, _ = os.Hostname()

	AppEnv = env.Get("APP_ENV", "local")

	IsLocal      = AppEnv == "local"
	IsProduction = AppEnv == "production"
	IsStaging    = AppEnv == "staging"
	IsDebug      = false
)

const (
	PgConnectCtxKey    = "PgConnectCtxKey"
	RedisConnectCtxKey = "RedisConnectCtxKey"
	AuthHeaderToken    = "AuthHeaderTokenCtxKey"
	TokenPayloadCtxKey = "TokenPayloadCtxKey"
	RequestIdCtxKey    = "RequestIdCtxKey"
	LoggerCtxKey       = "LoggerCtxKey"
)

func GetShutdownTimeout() time.Duration {
	timeout, err := strconv.Atoi(env.Get("SHUTDOWN_TIMEOUT", "0"))
	if err != nil {
		return 0
	}
	return time.Duration(timeout)
}
