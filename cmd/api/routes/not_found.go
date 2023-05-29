package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func notFound(c *gin.Context) {
	c.Error(errors.New(errors.Input{
		Message:    fmt.Sprintf("Cannot %s %s", c.Request.Method, c.Request.URL.String()),
		StatusCode: http.StatusNotFound,
	}))
}
