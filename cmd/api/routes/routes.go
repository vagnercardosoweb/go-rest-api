package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/cmd/api/middlewares"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func healthy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"date":      time.Now().UTC(),
		"ipAddress": c.RemoteIP(),
		"userAgent": c.Request.UserAgent(),
		"path":      fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
	})
}

func favicon(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func notAllowed(c *gin.Context) {
	c.Error(errors.New(errors.Input{
		Message:    fmt.Sprintf("Not allowed %s %s", c.Request.Method, c.Request.URL.Path),
		StatusCode: http.StatusMethodNotAllowed,
	}))
}

func notFound(c *gin.Context) {
	c.Error(errors.New(errors.Input{
		Message:    fmt.Sprintf("Cannot %s %s", c.Request.Method, c.Request.URL.String()),
		StatusCode: http.StatusNotFound,
	}))
}

func Setup(router *gin.Engine) {
	router.NoRoute(notFound)
	router.NoMethod(notAllowed)

	router.GET("/", middlewares.NoCacheHandler, healthy)
	router.GET("/favicon.ico", favicon)
}
