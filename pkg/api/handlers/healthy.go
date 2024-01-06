package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"net/http"
)

func Healthy(c *gin.Context) any {
	redisPingError := utils.GetRedisClient(c).Ping()
	pgPingError := utils.GetPgClient(c).Ping()

	if redisPingError != nil || pgPingError != nil {
		utils.GetLogger(c).
			AddMetadata("redisPingError", redisPingError).
			AddMetadata("pgPingError", pgPingError).
			Error("HEALTHY_ERROR")
		c.Writer.WriteHeader(http.StatusServiceUnavailable)
		return "UNAVAILABLE"
	}

	return "OK"
}
