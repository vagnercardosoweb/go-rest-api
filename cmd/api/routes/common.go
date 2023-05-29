package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func favicon(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func healthy(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"date":      time.Now().UTC(),
		"ipAddress": c.RemoteIP(),
		"userAgent": c.Request.UserAgent(),
		"path":      fmt.Sprintf("%s %s", c.Request.Method, c.Request.URL.Path),
	})
}
