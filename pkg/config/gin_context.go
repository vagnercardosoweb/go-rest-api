package config

import (
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

const (
	PgClientCtxKey      = "PgClientCtxKey"
	AuthTokenCtxKey     = "AuthTokenCtxKey"
	RedisClientCtxKey   = "RedisClientCtxKey"
	TokenDecodedCtxKey  = "TokenDecodedCtxKey"
	RequestLoggerCtxKey = "RequestLoggerCtxKey"
	TokenManagerCtxKey  = "TokenManagerCtxKey"
	RequestIdCtxKey     = "RequestIdCtxKey"
)

func TokenFromCtx(c *gin.Context) token.Token {
	return c.MustGet(TokenManagerCtxKey).(token.Token)
}

func AuthTokenFromCtx(c *gin.Context) string {
	return c.GetString(AuthTokenCtxKey)
}

func PgClientFromCtx(c *gin.Context) *postgres.Client {
	return c.MustGet(PgClientCtxKey).(*postgres.Client)
}

func RedisClientFromCtx(c *gin.Context) *redis.Client {
	return c.MustGet(RedisClientCtxKey).(*redis.Client)
}

func TokenDecodedFromCtx(c *gin.Context) *token.Decoded {
	return c.MustGet(TokenDecodedCtxKey).(*token.Decoded)
}

func LoggerFromCtx(c *gin.Context) *logger.Logger {
	return c.MustGet(RequestLoggerCtxKey).(*logger.Logger)
}

func RequestIdFromCtx(c *gin.Context) string {
	return c.GetString(RequestIdCtxKey)
}
