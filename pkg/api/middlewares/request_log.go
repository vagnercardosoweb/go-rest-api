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
}

func RequestLog(c *gin.Context) {
	path := c.Request.URL.Path

	if slices.Contains(skipPaths, path) {
		c.Next()
		return
	}

	method := c.Request.Method
	requestLogger := apicontext.Logger(c)
	logData := make(map[string]any)

	if routePath := c.FullPath(); routePath != "" {
		logData["routePath"] = routePath
	}

	logData["ip"] = c.ClientIP()
	logData["request"] = fmt.Sprintf("%s %s", method, path)
	logData["queryParams"] = c.Request.URL.Query()
	logData["userAgent"] = c.Request.UserAgent()
	logData["time"] = time.Since(apicontext.StartTime(c)).String()

	requestLogger.
		WithFields(logData).
		Info("HTTP_REQUEST_STARTED")

	// Process request
	c.Next()

	status := c.Writer.Status()

	logData["statusCode"] = status
	logData["time"] = c.Writer.Header().Get("X-Response-Time")

	level := logger.LevelInfo
	if status < http.StatusOK || status >= http.StatusBadRequest {
		level = logger.LevelError
	}

	requestLogger.
		WithFields(logData).
		Log(level, "HTTP_REQUEST_COMPLETED")
}
