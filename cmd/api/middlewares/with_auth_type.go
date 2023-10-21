package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func WithAuthType(authType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoded := config.TokenDecodedFromCtx(c)
		tokenType := decoded.Meta["type"].(string)

		if tokenType != authType {
			c.Error(errors.New(errors.Input{
				Name:        "UnauthorizedWithLogoutError",
				Message:     "Access token is not valid, please login again",
				StatusCode:  http.StatusUnauthorized,
				SendToSlack: errors.Bool(false),
				Metadata: errors.Metadata{
					"tokenType":    tokenType,
					"requiredType": authType,
				},
			}))
			c.Abort()
			return
		}

		c.Next()
	}
}
