package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/vagnercardosoweb/go-rest-api/cmd/api/middlewares"
)

func Setup(router *gin.Engine) {
	router.NoRoute(notFound)
	router.NoMethod(notAllowed)

	router.GET("/", middlewares.WrapHandler(healthy))
	router.GET("/favicon.ico", favicon)
}
