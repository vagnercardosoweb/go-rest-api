package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	apiresponse "github.com/vagnercardosoweb/go-rest-api/pkg/api/response"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func AuthType(authType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		decoded := apicontext.TokenOutput(c)
		tokenType := decoded.Meta["type"].(string)

		if tokenType != authType {
			apiresponse.Error(c, errors.New(errors.Input{
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
