package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func favicon(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusOK)
}

func healthy(*gin.Context) any {
	return "🚀"
}
