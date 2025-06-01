package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func Recovery(c *gin.Context, err any) {
	_ = c.Error(errors.New(errors.Input{
		Code:          "PANIC_ERROR",
		Message:       "The application received a panic error",
		OriginalError: err,
	}))

	ResponseError(c)
}
