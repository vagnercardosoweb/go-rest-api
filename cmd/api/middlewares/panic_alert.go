package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func panicAlert(c *gin.Context, err any) {
	message := err

	if e, ok := message.(error); ok {
		message = e.Error()
	}

	c.Error(errors.New(errors.Input{
		Code:          "PANIC_ERROR",
		Message:       "The application received a panic error",
		SendToSlack:   true,
		StatusCode:    http.StatusInternalServerError,
		OriginalError: message,
	}))

	responseError(c)
}
