package routes

import (
	"github.com/gin-gonic/gin"
)

func Setup(router *gin.Engine) {
	setupV1(router)
	setupCommon(router)
}
