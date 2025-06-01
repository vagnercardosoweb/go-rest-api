package apiresponse

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/utils"
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

	c.JSON(status, gin.H{
		"data":        data,
		"path":        fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.String()),
		"duration":    time.Since(apicontext.StartTime(c)).String(),
		"hostname":    utils.Hostname,
		"environment": env.GetAppEnv(),
		"requestId":   c.Writer.Header().Get("X-Request-Id"),
		"ipAddress":   c.ClientIP(),
		"userAgent":   c.Request.UserAgent(),
		"timezone":    time.UTC.String(),
		"brlDate":     utils.NowBrl(),
		"utcDate":     utils.NowUtc(),
	})
}
