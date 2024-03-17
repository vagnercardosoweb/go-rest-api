package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func NotFound(c *gin.Context) {
	if c.Request.URL.Path == "/" {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	utils.AbortWithError(c, errors.New(errors.Input{
		Message:    fmt.Sprintf("Cannot %s %s", c.Request.Method, c.Request.URL.String()),
		StatusCode: http.StatusNotFound,
		Logging:    errors.Bool(false),
	}))
}
