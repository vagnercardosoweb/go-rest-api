package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger_new"
)

func requestId(c *gin.Context) {
	requestId := c.GetHeader("x-amzn-trace-id")
	if requestId == "" {
		requestId = c.GetHeader("x-amzn-requestid")
	}
	if requestId == "" {
		requestId = uuid.New().String()
	}
	c.Set(config.RequestIdCtxKey, requestId)
	c.Set(config.LoggerCtxKey, logger_new.New().WithID(requestId))
	c.Header("X-Request-Id", requestId)
	c.Next()
}
