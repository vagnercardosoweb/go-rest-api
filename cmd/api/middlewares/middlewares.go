package middlewares

import "github.com/gin-gonic/gin"

func Setup(router *gin.Engine) {
	router.Use(cors)
	router.Use(responseTimer)
	router.Use(requestId)
	router.Use(extractAuthToken)
	router.Use(loggerRequest)
	router.Use(gin.CustomRecovery(panicAlert))
	router.Use(responseError)
}
