package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/cmd/api/handlers"
	"github.com/vagnercardosoweb/go-rest-api/cmd/api/middlewares"
)

func Setup(router *gin.Engine) {
	router.NoRoute(handlers.NotFound)
	router.NoMethod(handlers.NotAllowed)

	router.GET("/", middlewares.WrapHandler(handlers.Healthy))
	router.GET("/favicon.ico", handlers.Favicon)
}
