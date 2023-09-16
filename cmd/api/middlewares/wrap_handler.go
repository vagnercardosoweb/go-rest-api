package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
)

var (
	hostname, _ = os.Hostname()
)

func WrapHandler(handler func(c *gin.Context) any) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := handler(c)

		if err, ok := result.(error); ok {
			c.Error(err)
			c.Abort()
			return
		}

		status := c.Writer.Status()

		if result == nil && (status == http.StatusOK || status == 0) {
			c.Writer.WriteHeader(http.StatusNoContent)
			return
		}

		if result != nil {
			now := time.Now()
			c.JSON(status, gin.H{
				"data":        result,
				"path":        fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.String()),
				"duration":    time.Since(c.Writer.(*XResponseTimer).Start).String(),
				"hostname":    hostname,
				"environment": config.AppEnv,
				"ipAddress":   c.ClientIP(),
				"userAgent":   c.Request.UserAgent(),
				"timezone":    time.UTC.String(),
				"brlDate":     now.In(config.LocationBrl),
				"utcDate":     now.UTC(),
			})
		}

		c.Next()
	}
}
