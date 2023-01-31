package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/api/handlers"
)

func setupCommon(router *gin.Engine) {
	router.GET("/", handlers.Healthy)
	router.GET("/favicon.ico", handlers.Favicon)

	router.NoRoute(handlers.NotFound)
	router.NoMethod(handlers.NotAllowed)
}
