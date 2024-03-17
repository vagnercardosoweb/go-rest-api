package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func WithAuthType(authType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoded := utils.GetTokenOutput(c)
		tokenType := decoded.Meta["type"].(string)

		if tokenType != authType {
			utils.AbortWithError(c, errors.New(errors.Input{
				Name:        "UnauthorizedWithLogoutError",
				Message:     "Access token is not valid, please login again",
				StatusCode:  http.StatusUnauthorized,
				SendToSlack: errors.Bool(false),
				Metadata: errors.Metadata{
					"requiredType": authType,
					"tokenType":    tokenType,
				},
			}))
			return
		}

		c.Next()
	}
}
