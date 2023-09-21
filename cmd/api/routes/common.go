package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/cmd/api/middlewares"

	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func favicon(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func healthy(*gin.Context) any {
	return "🚀"
}

func notFound(c *gin.Context) {
	c.Error(errors.New(errors.Input{
		Message:     fmt.Sprintf("Cannot %s %s", c.Request.Method, c.Request.URL.String()),
		StatusCode:  http.StatusNotFound,
		SendToSlack: errors.Bool(false),
		Logging:     errors.Bool(false),
	}))
}

func notAllowed(c *gin.Context) {
	c.Error(errors.New(errors.Input{
		Message:     fmt.Sprintf("Not allowed %s %s", c.Request.Method, c.Request.URL.Path),
		StatusCode:  http.StatusMethodNotAllowed,
		SendToSlack: errors.Bool(false),
		Logging:     errors.Bool(false),
	}))
}

func makeCommon(router *gin.Engine) {
	router.NoRoute(notFound)
	router.NoMethod(notAllowed)

	router.GET("/", middlewares.WrapHandler(healthy))
	router.GET("/favicon.ico", favicon)
}
