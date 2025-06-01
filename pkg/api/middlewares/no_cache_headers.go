package middlewares

import (
	"github.com/gin-gonic/gin"
)

func NoCacheHeaders(c *gin.Context) {
	c.Header("Expires", "0")
	c.Header("Surrogate-Control", "no-store")
	c.Header("Pragma", "no-cache")
	c.Header(
		"Cache-Control",
		"no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0, post-check=0, pre-check=0",
	)
	c.Next()
}
