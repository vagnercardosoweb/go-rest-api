package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Favicon(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNoContent)
}
