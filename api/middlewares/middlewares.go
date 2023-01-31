package middlewares

import "github.com/gin-gonic/gin"

func Setup(router *gin.Engine) {
	router.Use(corsHandler)
	router.Use(requestIdHandler)
	router.Use(extractTokenHandler)
	router.Use(loggerHandler)
	router.Use(errorHandler)
}
