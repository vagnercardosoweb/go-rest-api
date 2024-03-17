package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
)

func ExtractAuthToken(c *gin.Context) {
	token := c.Query("token")
	authorization := c.GetHeader("Authorization")

	if authorization != "" {
		tokenParts := strings.Split(authorization, " ")

		if len(tokenParts) == 2 {
			token = tokenParts[1]
		}
	}

	c.Set(utils.AuthTokenCtxKey, token)
	c.Next()
}
