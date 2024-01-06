package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Favicon(c *gin.Context) {
	c.Writer.WriteHeader(http.StatusNoContent)
}
