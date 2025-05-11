package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
)

func Healthy(c *gin.Context) any {
	pgPingError := apicontext.PgClient(c).Ping()
	redisPingError := apicontext.RedisClient(c).Ping()

	if redisPingError != nil || pgPingError != nil {
		apicontext.Logger(c).
			AddMetadata("pgPingError", pgPingError).
			AddMetadata("redisPingError", redisPingError).
			Error("HEALTHY_ERROR")

		c.Writer.WriteHeader(http.StatusServiceUnavailable)

		return "UNAVAILABLE"
	}

	return "OK"
}
