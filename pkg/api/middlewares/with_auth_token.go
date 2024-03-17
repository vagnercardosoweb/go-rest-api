package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	tokenpkg "github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

var unauthorized = errors.New(errors.Input{
	Name:        "UnauthorizedWithLogoutError",
	StatusCode:  http.StatusUnauthorized,
	SendToSlack: errors.Bool(false),
	Message:     "Missing token in request.",
	Code:        "INVALID_JWT_TOKEN",
})

func WithAuthToken(c *gin.Context) {
	authToken := utils.GetAuthToken(c)
	if authToken == "" {
		utils.AbortWithError(c, unauthorized)
		return
	}

	if len(strings.Split(authToken, ".")) != 3 {
		unauthorized.Message = "The token is badly formatted."
		utils.AbortWithError(c, unauthorized)
		return
	}

	token := utils.GetTokenClient(c)
	unauthorized.Message = "Your access token is not valid, please login again."
	decoded, err := token.Decode(authToken)

	if err != nil {
		unauthorized.OriginalError = err.Error()
		utils.AbortWithError(c, unauthorized)
		return
	}

	if _, ok := decoded.Meta["type"]; !ok {
		utils.AbortWithError(c, unauthorized)
		return
	}

	c.Set(tokenpkg.OutputCtxKey, decoded)
	c.Next()
}
