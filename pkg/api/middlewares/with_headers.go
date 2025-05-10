package middlewares

import "github.com/gin-gonic/gin"

func WithHeaders(c *gin.Context) {
	// Not Expose Server Information
	c.Header("Server", "Go-API")

	// Cache Headers
	c.Header("Expires", "0")
	c.Header("Surrogate-Control", "no-store")
	c.Header("Pragma", "no-cache")
	c.Header(
		"Cache-Control",
		"no-store, no-cache, must-revalidate, proxy-revalidate, max-age=0, post-check=0, pre-check=0",
	)

	// Response Headers
	c.Header("Content-Type", "application/json; charset=utf-8")

	// Protection Headers
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
	c.Header("Content-Security-Policy", "default-src 'self'")
	c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")
	c.Header("Cross-Origin-Embedder-Policy", "require-corp")
	c.Header("Cross-Origin-Resource-Policy", "same-origin")
	c.Header("Cross-Origin-Opener-Policy", "same-origin")
	c.Header("X-Content-Type-Options", "nosniff")

	c.Next()
}
