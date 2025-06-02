package apiresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}

func Wrapper(handler func(c *gin.Context) any) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := handler(c)

		if err, ok := result.(error); ok {
			Error(c, err)
			return
		}

		Json(c, result)
	}
}

func Json(c *gin.Context, data any) {
	status := c.Writer.Status()

	if data == nil && (status == http.StatusOK || status == 0) {
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	if data == nil {
		return
	}

	c.JSON(status, data)
}
