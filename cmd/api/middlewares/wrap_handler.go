package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"net/http"
	"time"
)

func WrapHandler(handler func(c *gin.Context) any) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := handler(c)

		if result != nil {
			if err, ok := result.(*errors.Input); ok {
				c.AbortWithError(err.StatusCode, err)
				return
			} else if err, ok := result.(error); ok {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
		}

		status := c.Writer.Status()
		if status == 0 {
			status = http.StatusOK
		}

		c.JSON(status, gin.H{
			"data":      result,
			"path":      fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.String()),
			"ipAddress": c.ClientIP(),
			"timestamp": time.Now().UTC(),
			"duration":  time.Since(c.Writer.(*XResponseTimer).start).String(),
		})
		return
	}
}
