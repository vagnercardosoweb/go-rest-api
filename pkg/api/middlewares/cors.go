package middlewares

import (
	"net/http"
	"strings"

	"slices"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

func Cors(c *gin.Context) {
	origin := c.Request.Header.Get("Origin")
	allowedOrigins := env.GetAsString("CORS_ALLOW_ORIGINS", "*")

	if allowedOrigins == "*" {
		c.Header("Access-Control-Allow-Origin", "*")
	} else if isOriginAllowed(origin, strings.Split(allowedOrigins, ",")) {
		c.Header("Access-Control-Allow-Origin", origin)
	}

	c.Header("Access-Control-Allow-Methods", getAllowedMethods())
	c.Header("Access-Control-Max-Age", env.GetAsString("CORS_MAX_AGE", "84600"))
	c.Header("Access-Control-Allow-Credentials", env.GetAsString("CORS_ALLOW_CREDENTIALS", "true"))
	c.Header("Access-Control-Allow-Headers", getAllowedHeaders())

	if c.Request.Method == http.MethodOptions {
		c.Header("Content-Length", "0")
		c.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	c.Next()
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	return slices.Contains(allowedOrigins, origin)
}

func getAllowedMethods() string {
	return env.GetAsString("CORS_ALLOW_METHODS", "GET, POST, PUT, DELETE, PATCH, HEAD")
}

func getAllowedHeaders() string {
	return env.GetAsString(
		"CORS_ALLOW_HEADERS",
		strings.Join([]string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"Origin",
			"Cache-Control",
			"X-Requested-With",
			"X-HTTP-Method-Override",
			"Accept-Language",
			"X-Refresh-Token",
			"X-Id-Token",
		}, ", "),
	)
}
