package handlers

import (
	"net/http"
	"rest-api/errors"

	"github.com/gin-gonic/gin"
)

func NotAllowed(ctx *gin.Context) {
	notAllowedError := errors.NewMethodNotAllowed(
		"Method not allowed",
		map[string]interface{}{
			"path":   ctx.Request.URL.Path,
			"method": ctx.Request.Method,
		})
	ctx.JSON(http.StatusMethodNotAllowed, notAllowedError)
	return
}
