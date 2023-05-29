package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func notAllowed(c *gin.Context) {
	c.Error(errors.New(errors.Input{
		Message:    fmt.Sprintf("Not allowed %s %s", c.Request.Method, c.Request.URL.Path),
		StatusCode: http.StatusMethodNotAllowed,
	}))
}
