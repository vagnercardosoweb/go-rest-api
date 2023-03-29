package middlewares

import "github.com/gin-gonic/gin"

func NoCacheHandler(c *gin.Context) {
	c.Header("Expires", "0")
	c.Header("Pragma", "no-cache")
	c.Header("Surrogate-Control", "no-store")
	c.Header(
		"Cache-Control",
		"no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0, post-check=0, pre-check=0",
	)
	c.Next()
}
