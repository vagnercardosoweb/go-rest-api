package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

func SecurityHeaders(c *gin.Context) {
	// Remove headers that can leak information
	c.Header("Server", "")
	c.Header("X-Powered-By", "")

	// Protection Headers
	c.Header("X-Frame-Options", "DENY")
	c.Header("X-XSS-Protection", "1; mode=block")
	c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

	if env.IsProduction() {
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
	}

	c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=(), payment=()")
	c.Header("Cross-Origin-Embedder-Policy", "require-corp")
	c.Header("Cross-Origin-Resource-Policy", "same-origin")
	c.Header("Cross-Origin-Opener-Policy", "same-origin")
	c.Header("X-Content-Type-Options", "nosniff")

	// Content Security Policy
	c.Header(
		"Content-Security-Policy",
		"default-src 'self'; "+
			"script-src 'self' 'unsafe-inline'; "+
			"style-src 'self' 'unsafe-inline'; "+
			"img-src 'self' data: https:; "+
			"font-src 'self'; "+
			"connect-src 'self'; "+
			"frame-ancestors 'none'",
	)

	c.Next()
}
