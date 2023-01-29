package middlewares

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
)

func requestIdHandler(c *gin.Context) {
	requestId := c.GetHeader("x-amzn-trace-id")
	if requestId == "" {
		requestId = uuid.NewV4().String()
	}
	c.Set(config.RequestIdContextKey, requestId)
	c.Header("X-Request-Id", requestId)
	c.Next()
}
