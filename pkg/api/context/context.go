package apicontext

import (
	"time"

	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

const (
	StartTimeKey           = "StartTimeKey"
	BearerTokenKey         = "BearerTokenKey"
	ValidatorTranslatorKey = "ValidatorTranslatorKey"
	RequestIdKey           = "RequestIdKey"
)

func PgClient(c *gin.Context) *postgres.Client {
	return c.MustGet(postgres.CtxKey).(*postgres.Client)
}

func TokenClient(c *gin.Context) token.Client {
	return c.MustGet(token.CtxClientKey).(token.Client)
}

func TokenOutput(c *gin.Context) *token.Output {
	return c.MustGet(token.CtxDecodedKey).(*token.Output)
}

func PasswordHasher(c *gin.Context) password.PasswordHasher {
	return c.MustGet(password.CtxKey).(password.PasswordHasher)
}

func BearerToken(c *gin.Context) string {
	return c.GetString(BearerTokenKey)
}

func RedisClient(c *gin.Context) *redis.Client {
	return c.MustGet(redis.CtxKey).(*redis.Client)
}

func ValidatorTranslator(c *gin.Context) *ut.Translator {
	return c.MustGet(ValidatorTranslatorKey).(*ut.Translator)
}

func Logger(c *gin.Context) *logger.Logger {
	return c.MustGet(logger.CtxKey).(*logger.Logger)
}

func RequestId(c *gin.Context) string {
	return c.GetString(RequestIdKey)
}

func StartTime(c *gin.Context) time.Time {
	return c.MustGet(StartTimeKey).(time.Time)
}
