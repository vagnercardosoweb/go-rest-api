package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	origin  = "*"
	methods = "GET, POST, PUT, DELETE, PATCH, OPTIONS"
	headers = "Accept, Origin, Content-Type, Authorization, Cache-Control, X-Requested-With, X-HTTP-Method-Override, Accept-Language, X-Refresh-Token, X-Id-Token, X-Aws-IdToken"
)

func corsHandler(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", origin)
	c.Header("Access-Control-Allow-Methods", methods)
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers", headers)

	if c.Request.Method == http.MethodOptions {
		c.Writer.WriteHeader(http.StatusOK)
	}

	c.Next()
}
