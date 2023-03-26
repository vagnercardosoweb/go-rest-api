package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

func loggerHandler(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	routePath := c.FullPath()
	method := c.Request.Method
	requestId := c.GetString(config.RequestIdCtxKey)
	loggerId := fmt.Sprintf("REQ:%s", requestId)
	clientIP := c.ClientIP()

	metadata := logger.Metadata{
		"ip":        clientIP,
		"path":      path,
		"full_path": routePath,
		"method":    method,
		"query":     c.Request.URL.Query(),
		"version":   c.Request.Proto,
		"referer":   c.GetHeader("referer"),
		"agent":     c.Request.UserAgent(),
		"time":      0,
		"length":    0,
		"status":    0,
	}

	logger.Log(logger.Input{
		Id:       loggerId,
		Level:    logger.DEBUG,
		Message:  "STARTED",
		Metadata: metadata,
	})

	// Process request
	c.Next()

	end := time.Now()
	status := c.Writer.Status()
	latency := end.Sub(start)

	metadata["time"] = latency.String()
	metadata["length"] = c.Writer.Size()
	metadata["status"] = status

	if config.IsDebug && method != http.MethodGet {
		metadata["body"] = c.Request.Form
	}

	level := logger.DEBUG
	if status >= http.StatusInternalServerError {
		level = logger.ERROR
	}

	logger.Log(logger.Input{
		Id:       loggerId,
		Level:    level,
		Message:  "COMPLETED",
		Metadata: metadata,
	})
}
