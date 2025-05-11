package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/vagnercardosoweb/go-rest-api/pkg/api/response"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func NotAllowed(c *gin.Context) {
	apiresponse.Error(c, errors.New(errors.Input{
		Message:    fmt.Sprintf("Not allowed %s %s", c.Request.Method, c.Request.URL.Path),
		StatusCode: http.StatusMethodNotAllowed,
		Logging:    errors.Bool(false),
	}))
}
