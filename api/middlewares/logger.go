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
	logger := logger.New(logger.Input{Id: fmt.Sprintf("REQ:%s", requestId)})

	// Set logger context
	c.Set(config.LoggerContextKey, logger)

	logger.
		AddMetadata("ip", c.ClientIP()).
		AddMetadata("path", path).
		AddMetadata("route_path", routePath).
		AddMetadata("method", method).
		AddMetadata("query", c.Request.URL.Query()).
		AddMetadata("version", c.Request.Proto).
		AddMetadata("referer", c.GetHeader("referer")).
		AddMetadata("agent", c.Request.UserAgent()).
		Info("started")

	// Process request
	c.Next()

	end := time.Now()
	latency := end.Sub(start)
	status := c.Writer.Status()

	logger.
		AddMetadata("time", latency.String()).
		AddMetadata("time_ms", latency.Milliseconds()).
		AddMetadata("status", status).
		AddMetadata("length", c.Writer.Size())

	if config.IsDebug && c.Request.Method != http.MethodGet {
		logger.AddMetadata("body", c.Request.Form)
	}

	logLevel := "INFO"
	if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
		logLevel = "WARN"
	} else if status >= http.StatusInternalServerError {
		logLevel = "ERROR"
	}

	logger.Log(logLevel, "completed")
}
