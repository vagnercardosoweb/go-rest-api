package utils

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	ut "github.com/go-playground/universal-translator"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
	"github.com/vagnercardosoweb/go-rest-api/pkg/utils"
	"io"
	"net/http"
	"time"
)

func AbortWithError(c *gin.Context, err error) {
	c.Error(err)
	c.Abort()
}

func ResponseWithSuccess(c *gin.Context, data any) {
	status := c.Writer.Status()

	if data == nil && (status == http.StatusOK || status == 0) {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	c.JSON(status, gin.H{
		"data":        data,
		"path":        fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.String()),
		"duration":    time.Since(GetRequestStartTime(c)).String(),
		"hostname":    utils.Hostname,
		"environment": env.GetAppEnv(),
		"requestId":   c.Writer.Header().Get("X-Request-Id"),
		"ipAddress":   c.ClientIP(),
		"userAgent":   c.Request.UserAgent(),
		"timezone":    time.UTC.String(),
		"brlDate":     time.Now().In(utils.LocationBrl),
		"utcDate":     time.Now().UTC(),
	})
}

func WrapperHandler(handler func(c *gin.Context) any) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := handler(c)

		if err, ok := result.(error); ok {
			AbortWithError(c, err)
			return
		}

		ResponseWithSuccess(c, result)
	}
}

func GetBodyAsBytes(c *gin.Context) []byte {
	bodyAsBytes := []byte("{}")
	if val, ok := c.Get(gin.BodyBytesKey); ok && val != nil {
		bodyAsBytes = val.([]byte)
	} else {
		b, _ := io.ReadAll(c.Request.Body)
		if len(b) > 0 {
			c.Set(gin.BodyBytesKey, b)
			bodyAsBytes = b
		}
	}
	return bodyAsBytes
}

func GetBodyAsJson(c *gin.Context) map[string]any {
	bodyAsBytes := GetBodyAsBytes(c)
	result := make(map[string]any)
	_ = json.Unmarshal(bodyAsBytes, &result)
	return result
}

func GetTokenClient(c *gin.Context) token.Client {
	return c.MustGet(TokenClientCtxKey).(token.Client)
}

func GetAuthToken(c *gin.Context) string {
	return c.GetString(AuthTokenCtxKey)
}

func GetPgClient(c *gin.Context) *postgres.Client {
	return c.MustGet(PgClientCtxKey).(*postgres.Client)
}

func GetRedisClient(c *gin.Context) *redis.Client {
	return c.MustGet(RedisClientCtxKey).(*redis.Client)
}

func GetTokenDecoded(c *gin.Context) *token.Output {
	return c.MustGet(TokenDecodedCtxKey).(*token.Output)
}

func GetLogger(c *gin.Context) *logger.Logger {
	return c.MustGet(RequestLoggerCtxKey).(*logger.Logger)
}

func GetRequestId(c *gin.Context) string {
	return c.GetString(RequestIdCtxKey)
}

func GetRequestStartTime(c *gin.Context) time.Time {
	return c.MustGet(RequestStartTimeKey).(time.Time)
}

func GetValidateTranslator(c *gin.Context) ut.Translator {
	return c.MustGet(ValidateTranslatorCtxKey).(ut.Translator)
}
