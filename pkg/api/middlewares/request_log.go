package middlewares

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

var skipPaths = []string{
	"/healthy",
	"/favicon.ico",
	"/timestamp",
	"/",
}

func RequestLog(c *gin.Context) {
	path := c.Request.URL.Path

	if slices.Contains(skipPaths, path) {
		c.Next()
		return
	}

	method := c.Request.Method
	requestLogger := apicontext.Logger(c)
	clientIP := c.ClientIP()

	metadata := map[string]any{
		"ip":        clientIP,
		"request":   fmt.Sprintf("%s %s", method, path),
		"userAgent": c.Request.UserAgent(),
		"time":      time.Since(apicontext.StartTime(c)).String(),
	}

	requestLogger.
		WithoutRedact().
		WithMetadata(metadata).
		Info("HTTP_REQUEST_STARTED")

	// Process request
	c.Next()

	status := c.Writer.Status()

	metadata["statusCode"] = status
	metadata["time"] = c.Writer.Header().Get("X-Response-Time")

	level := logger.LevelInfo
	if status < http.StatusOK || status >= http.StatusBadRequest {
		level = logger.LevelError
	}

	requestLogger.
		WithoutRedact().
		WithMetadata(metadata).
		Log(level, "HTTP_REQUEST_COMPLETED")
}
