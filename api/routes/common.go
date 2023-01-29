package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/api/handlers"
)

func setupCommon(router *gin.Engine) {
	router.NoRoute(handlers.NotFound)
	router.NoMethod(handlers.NotAllowed)

	router.GET("/", handlers.Healthy)
}
