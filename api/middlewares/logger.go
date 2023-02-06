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
	requestId := c.GetString(config.RequestIdContextKey)
	levelId := fmt.Sprintf("REQ:%s", requestId)
	logger.Log(logger.Input{
		Id:    levelId,
		Level: logger.DEBUG,
		Metadata: logger.Metadata{
			"ip":         c.ClientIP(),
			"path":       path,
			"route_path": routePath,
			"method":     method,
			"query":      c.Request.URL.Query(),
			"version":    c.Request.Proto,
			"referer":    c.GetHeader("referer"),
			"agent":      c.Request.UserAgent(),
		},
	})

	// Process request
	c.Next()

	end := time.Now()
	status := c.Writer.Status()
	latency := end.Sub(start)

	metadata := logger.Metadata{
		"time":   latency.String(),
		"length": c.Writer.Size(),
		"status": status,
	}

	if config.IsDebug && method != http.MethodGet {
		metadata["body"] = c.Request.Form
	}

	level := logger.DEBUG
	if status >= http.StatusInternalServerError {
		level = logger.ERROR
	}

	logger.Log(logger.Input{
		Id:       levelId,
		Level:    level,
		Message:  "completed",
		Metadata: metadata,
	})
}
