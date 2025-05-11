package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/vagnercardosoweb/go-rest-api/pkg/api/response"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func NotFound(c *gin.Context) {
	if c.Request.URL.Path == "/" {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	apiresponse.Error(c, errors.New(errors.Input{
		Message:    fmt.Sprintf("Cannot %s %s", c.Request.Method, c.Request.URL.String()),
		StatusCode: http.StatusNotFound,
		Logging:    errors.Bool(false),
	}))
}
