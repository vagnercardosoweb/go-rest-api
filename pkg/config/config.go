package config

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

const (
	TokenCtxKey         = "TokenCtxKey"
	PgClientCtxKey      = "PgClientCtxKey"
	RedisClientCtxKey   = "RedisClientCtxKey"
	AuthHeaderTokenCtx  = "AuthHeaderTokenCtxKey"
	TokenDecodedCtxKey  = "TokenDecodedCtxKey"
	RequestLoggerCtxKey = "RequestLoggerCtxKey"
	RequestIdCtxKey     = "RequestIdCtxKey"
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

func GetAppEnv() string {
	return env.Get("APP_ENV", "local")
}

func GetPid() int {
	return os.Getpid()
}

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func GetShutdownTimeout() time.Duration {
	timeout, err := strconv.Atoi(env.Get("SHUTDOWN_TIMEOUT", "0"))
	if err != nil {
		return 0
	}
	return time.Duration(timeout)
}

func GetLocationUtc() *time.Location {
	return time.UTC
}

func GetLocationGlobal() *time.Location {
	loc, _ := time.LoadLocation(env.Get("TZ", "UTC"))
	return loc
}

func GetLocationBrl() *time.Location {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	return loc
}

func GetExpiresInJwt() time.Duration {
	expiresInSecond, _ := strconv.Atoi(env.Get("JWT_EXPIRES_IN_SECONDS", "86400"))
	return time.Duration(expiresInSecond)
}

func GetTokenFromCtx(c *gin.Context) token.Token {
	return c.MustGet(TokenCtxKey).(token.Token)
}

func GetAuthHeaderTokenFromCtx(c *gin.Context) string {
	return c.GetString(AuthHeaderTokenCtx)
}

func GetPgClient(c *gin.Context) *postgres.Client {
	return c.MustGet(PgClientCtxKey).(*postgres.Client)
}

func GetRedisConnectFromCtx(c *gin.Context) *redis.Client {
	return c.MustGet(RedisClientCtxKey).(*redis.Client)
}

func GetTokenDecodedFromCtx(c *gin.Context) *token.Decoded {
	return c.MustGet(TokenDecodedCtxKey).(*token.Decoded)
}

func GetLoggerFromCtx(c *gin.Context) *logger.Logger {
	return c.MustGet(RequestLoggerCtxKey).(*logger.Logger)
}

func GetRequestIdFromCtx(c *gin.Context) string {
	return c.GetString(RequestIdCtxKey)
}
