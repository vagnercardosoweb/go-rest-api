package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/api/handlers"
)

func setupV1(router *gin.Engine) {
	v1 := router.Group("/v1")
	v1.GET("", handlers.Healthy)
}
