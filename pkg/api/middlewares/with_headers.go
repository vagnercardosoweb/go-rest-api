package middlewares

import "github.com/gin-gonic/gin"

func WithHeaders(c *gin.Context) {
	c.Header("Expires", "0")
	c.Header("Pragma", "no-cache")
	c.Header("Surrogate-Control", "no-store")
	c.Header(
		"Cache-Control",
		"no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0, post-check=0, pre-check=0",
	)
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Content-Security-Policy", "frame-ancestors 'none'")
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "0")
	c.Header("Referrer-Policy", "same-origin")
	c.Header("X-Frame-Options", "DENY")
	c.Next()
}
