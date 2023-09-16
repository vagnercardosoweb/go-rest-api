package config

import (
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

func GetTokenFromCtx(c *gin.Context) token.Token {
	return c.MustGet(TokenCtxKey).(token.Token)
}

func GetAuthHeaderTokenFromCtx(c *gin.Context) string {
	return c.GetString(AuthHeaderTokenCtx)
}

func GetPgClientFromCtx(c *gin.Context) *postgres.Client {
	return c.MustGet(PgClientCtxKey).(*postgres.Client)
}

func GetRedisClientFromCtx(c *gin.Context) *redis.Client {
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
