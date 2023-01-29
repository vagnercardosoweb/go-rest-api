package middlewares

import "github.com/gin-gonic/gin"

func noCacheHandler(c *gin.Context) {
	c.Header("Expires", "0")
	c.Header("Pragma", "no-cache")
	c.Header("Surrogate-Control", "no-store")
	c.Header(
		"Cache-Control",
		"no-store, no-cache, must-revalidate, proxy-revalidate",
	)
	c.Next()
}
