package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"net/http"
)

func NotAllowed(c *gin.Context) {
	utils.AbortWithError(c, errors.New(errors.Input{
		Message:    fmt.Sprintf("Not allowed %s %s", c.Request.Method, c.Request.URL.Path),
		StatusCode: http.StatusMethodNotAllowed,
		Logging:    errors.Bool(false),
	}))
}
