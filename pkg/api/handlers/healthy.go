package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/utils"
)

func Healthy(c *gin.Context) {
	pgPingError := apicontext.PgClient(c).Ping()
	redisPingError := apicontext.RedisClient(c).Ping()

	data := "OK"
	status := http.StatusOK

	if redisPingError != nil || pgPingError != nil {
		apicontext.Logger(c).
			AddField("pgPingError", pgPingError).
			AddField("redisPingError", redisPingError).
			Error("HEALTHY_ERROR")

		status = http.StatusServiceUnavailable
		data = "UNAVAILABLE"
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
