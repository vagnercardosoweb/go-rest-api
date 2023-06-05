package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"net/http"
	"time"
)

func WrapHandler(handler func(c *gin.Context) any) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := handler(c)

		if err, ok := result.(error); ok {
			c.Error(err)
			return
		}

		status := c.Writer.Status()

		if result == nil && (status == http.StatusOK || status == 0) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}

		if result != nil {
			c.JSON(status, gin.H{
				"data":      result,
				"path":      fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.String()),
				"ipAddress": c.ClientIP(),
				"timestamp": time.Now().UTC(),
				"duration":  time.Since(c.Writer.(*XResponseTimer).start).String(),
				"hostname":  config.Hostname,
			})
		}

		return
	}
}
