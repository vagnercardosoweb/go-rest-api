package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func Healthy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"date":      time.Now().UTC(),
		"hostname":  config.Hostname,
		"ipAddress": c.RemoteIP(),
		"userAgent": c.Request.UserAgent(),
	})
	return
}

func Favicon(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
	return
}

func NotAllowed(ctx *gin.Context) {
	notAllowedError := errors.NewMethodNotAllowed(errors.Input{
		Message: "Method not allowed",
		Metadata: errors.Metadata{
			"path":   ctx.Request.URL.Path,
			"method": ctx.Request.Method,
		},
	})
	ctx.JSON(http.StatusMethodNotAllowed, notAllowedError)
	return
}

func NotFound(ctx *gin.Context) {
	notFoundError := errors.NewNotFound(errors.Input{
		Message: "Page not found",
		Metadata: errors.Metadata{
			"path":   ctx.Request.URL.Path,
			"method": ctx.Request.Method,
		},
	})
	ctx.JSON(notFoundError.StatusCode, notFoundError)
	return
}
