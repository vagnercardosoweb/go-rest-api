package middlewares

import "github.com/gin-gonic/gin"

func Setup(router *gin.Engine) {
	router.Use(Cors)
	router.Use(ResponseTime)
	router.Use(RequestId)
	router.Use(ExtractToken)
	router.Use(RequestLog)
	router.Use(gin.CustomRecovery(PanicAlert))
	router.Use(ResponseError)
}
