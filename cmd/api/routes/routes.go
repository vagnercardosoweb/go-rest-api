package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/cmd/api/middlewares"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func healthy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"date":      time.Now().UTC(),
		"hostname":  config.Hostname,
		"ipAddress": c.RemoteIP(),
		"userAgent": c.Request.UserAgent(),
	})
}

func favicon(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func notAllowed(ctx *gin.Context) {
	notAllowedError := errors.NewMethodNotAllowed(errors.Input{
		Message: "Method not allowed",
		Metadata: errors.Metadata{
			"path":   ctx.Request.URL.Path,
			"method": ctx.Request.Method,
		},
	})
	ctx.JSON(http.StatusMethodNotAllowed, notAllowedError)
}

func notFound(ctx *gin.Context) {
	notFoundError := errors.NewNotFound(errors.Input{
		Message: "Page not found",
		Metadata: errors.Metadata{
			"path":   ctx.Request.URL.Path,
			"method": ctx.Request.Method,
		},
	})
	ctx.JSON(notFoundError.StatusCode, notFoundError)
}

func Setup(router *gin.Engine) {
	router.NoRoute(notFound)
	router.NoMethod(notAllowed)

	router.GET("/", middlewares.NoCacheHandler, healthy)
	router.GET("/favicon.ico", favicon)
}
