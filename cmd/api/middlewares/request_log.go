package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

var skipPaths = []string{
	"/", "/favicon.ico",
}

func RequestLog(c *gin.Context) {
	path := c.Request.URL.String()

	for _, skipPath := range skipPaths {
		if path == skipPath {
			c.Next()
			return
		}
	}

	method := c.Request.Method
	requestLogger := config.GetLoggerFromCtx(c)
	clientIP := c.ClientIP()
	metadata := map[string]any{
		"ip":      clientIP,
		"method":  method,
		"path":    path,
		"query":   c.Request.URL.Query(),
		"version": c.Request.Proto,
		"referer": c.GetHeader("referer"),
		"agent":   c.Request.UserAgent(),
		"time":    0,
		"length":  0,
		"status":  0,
	}

	if routePath := c.FullPath(); routePath != "" {
		metadata["route_path"] = routePath
	}

	requestLogger.WithMetadata(metadata).Info("HTTP_REQUEST_STARTED")

	// Process request
	c.Next()

	status := c.Writer.Status()

	metadata["time"] = c.Writer.Header().Get("X-Response-Time")
	metadata["length"] = c.Writer.Size()
	metadata["status"] = status

	if config.IsDebug && method != http.MethodGet {
		metadata["body"] = GetBodyAsJson(c)
	}

	level := logger.LevelInfo
	if status < http.StatusOK || status >= http.StatusBadRequest {
		level = logger.LevelError
	}

	requestLogger.
		WithMetadata(metadata).
		Log(
			level,
			"HTTP_REQUEST_COMPLETED",
		)
}
