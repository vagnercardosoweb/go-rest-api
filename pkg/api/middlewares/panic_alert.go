package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func Recovery(c *gin.Context, err any) {
	c.Error(errors.New(errors.Input{
		Code:          "PANIC_ERROR",
		Message:       "The application received a panic error",
		SendToSlack:   errors.Bool(true),
		StatusCode:    http.StatusInternalServerError,
		OriginalError: err,
	}))

	ResponseError(c)
}
