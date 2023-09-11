package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func Favicon(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func Healthy(*gin.Context) any {
	return "🚀"
}

func NotFound(c *gin.Context) {
	c.Error(errors.New(errors.Input{
		Message:     fmt.Sprintf("Cannot %s %s", c.Request.Method, c.Request.URL.String()),
		StatusCode:  http.StatusNotFound,
		SendToSlack: errors.Bool(false),
		Logging:     errors.Bool(false),
	}))
}

func NotAllowed(c *gin.Context) {
	c.Error(errors.New(errors.Input{
		Message:     fmt.Sprintf("Not allowed %s %s", c.Request.Method, c.Request.URL.Path),
		StatusCode:  http.StatusMethodNotAllowed,
		SendToSlack: errors.Bool(false),
		Logging:     errors.Bool(false),
	}))
}
