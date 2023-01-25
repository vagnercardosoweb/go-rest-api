package handlers

import (
	"net/http"
	"rest-api/config"
	"time"

	"github.com/gin-gonic/gin"
)

func Healthy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"date":      time.Now().UTC(),
		"hostname":  config.Hostname,
		"ipAddress": c.RemoteIP(),
		"userAgent": c.Request.UserAgent(),
	})
	return
}
