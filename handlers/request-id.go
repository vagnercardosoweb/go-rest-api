package handlers

import (
	"rest-api/config"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func RequestId(c *gin.Context) {
	requestId := c.GetHeader("x-amzn-trace-id")
	if requestId == "" {
		requestId = uuid.NewV4().String()
	}
	c.Set(config.RequestIdContextKey, requestId)
	c.Header("X-Request-Id", requestId)
	c.Next()
}
