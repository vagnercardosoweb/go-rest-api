package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
)

func BearerToken(c *gin.Context) {
	token := ""
	authorization := c.GetHeader("Authorization")

	if authorization != "" {
		tokenParts := strings.Fields(authorization)

		if len(tokenParts) == 2 && strings.EqualFold(tokenParts[0], "Bearer") {
			token = tokenParts[1]
		}
	}

	c.Set(apicontext.BearerTokenKey, strings.TrimSpace(token))
	c.Next()
}
