package handlers

import (
	"rest-api/errors"

	"github.com/gin-gonic/gin"
)

func NotFound(ctx *gin.Context) {
	notFoundError := errors.NewNotFound(
		"Page not found",
		map[string]interface{}{
			"path":   ctx.Request.URL.Path,
			"method": ctx.Request.Method,
		})
	ctx.JSON(notFoundError.StatusCode, notFoundError)
	return
}
