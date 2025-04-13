package middlewares

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

var skipPaths = []string{"/", "/healthy", "/timestamp", "/favicon.ico"}

func RequestLog(c *gin.Context) {
	path := c.Request.URL.Path

	if slices.Contains(skipPaths, path) {
		c.Next()
		return
	}

	method := c.Request.Method
	requestLogger := utils.GetLogger(c)
	clientIP := c.ClientIP()
	metadata := map[string]any{
		"ip":      clientIP,
		"method":  method,
		"path":    path,
		"query":   c.Request.URL.Query(),
		"version": c.Request.Proto,
		"referer": c.GetHeader("referer"),
		"agent":   c.Request.UserAgent(),
		"time":    time.Since(utils.GetRequestStartTime(c)).String(),
		"length":  0,
		"status":  0,
	}

	if routePath := c.FullPath(); routePath != "" {
		metadata["routePath"] = routePath
	}

	requestLogger.
		WithoutRedact().
		WithMetadata(metadata).
		Info("HTTP_REQUEST_STARTED")

	// Process request
	c.Next()

	status := c.Writer.Status()

	metadata["time"] = c.Writer.Header().Get("X-Response-Time")
	metadata["length"] = c.Writer.Size()
	metadata["status"] = status

	if method != http.MethodGet {
		metadata["body"] = utils.GetBodyAsJson(c)
	}

	level := logger.LevelInfo
	if status < http.StatusOK || status >= http.StatusBadRequest {
		level = logger.LevelError
	}

	requestLogger.
		WithMetadata(metadata).
		Log(level, "HTTP_REQUEST_COMPLETED")
}
