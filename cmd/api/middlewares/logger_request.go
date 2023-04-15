package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"

	"github.com/gin-gonic/gin"
)

func loggerRequest(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.String()
	routePath := c.FullPath()
	method := c.Request.Method
	requestId := c.GetString(config.RequestIdCtxKey)
	loggerId := fmt.Sprintf("REQ:%s", requestId)
	clientIP := c.ClientIP()

	metadata := logger.Metadata{
		"ip":       clientIP,
		"method":   method,
		"path":     path,
		"query":    c.Request.URL.Query(),
		"version":  c.Request.Proto,
		"referrer": c.GetHeader("referer"),
		"agent":    c.Request.UserAgent(),
		"time":     0,
		"length":   0,
		"status":   0,
	}

	if routePath != "" {
		metadata["route_path"] = routePath
	}

	logger.Log(logger.Input{
		Id:       loggerId,
		Level:    logger.DEBUG,
		Message:  "REQUEST_STARTED",
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
	if status < http.StatusOK || status >= http.StatusBadRequest {
		level = logger.ERROR
	}

	logger.Log(logger.Input{
		Id:       loggerId,
		Level:    level,
		Message:  "REQUEST_COMPLETED",
		Metadata: metadata,
	})
}
