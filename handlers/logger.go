package handlers

import (
	"fmt"
	"net/http"
	"time"

	"rest-api/config"
	"rest-api/shared"

	"github.com/gin-gonic/gin"
)

func Logger(c *gin.Context) {
	start := time.Now()
	path := c.Request.URL.Path
	method := c.Request.Method
	requestId := c.GetString(config.RequestIdContextKey)
	logger := shared.NewLogger(shared.Logger{Id: fmt.Sprintf("REQ:%s", requestId)})

	// Set logger context
	c.Set(config.LoggerContextKey, logger)

	logger.
		AddMetadata("ip", c.ClientIP()).
		AddMetadata("method", method).
		AddMetadata("path", path).
		AddMetadata("version", c.Request.Proto).
		AddMetadata("referer", c.GetHeader("referer")).
		AddMetadata("agent", c.Request.UserAgent())

	if config.IsDebug {
		if forwardedUser := c.GetHeader("X-Forwarded-User"); forwardedUser != "" {
			logger.AddMetadata("forwarded_user", forwardedUser)
		}
		if forwardedEmail := c.GetHeader("X-Forwarded-Email"); forwardedEmail != "" {
			logger.AddMetadata("forwarded_email", forwardedEmail)
		}
	}

	logger.Info("started")

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

	if config.IsDebug {
		logger.AddMetadata("raw_query", c.Request.URL.RawQuery)

		if c.Request.Method != http.MethodGet {
			logger.AddMetadata("raw_form", c.Request.Form)
		}
	}

	logLevel := "INFO"
	if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
		logLevel = "WARN"
	} else if status >= http.StatusInternalServerError {
		logLevel = "ERROR"
	}

	logger.Log(logLevel, "completed")
}
