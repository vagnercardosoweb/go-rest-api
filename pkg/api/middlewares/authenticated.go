package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	apiresponse "github.com/vagnercardosoweb/go-rest-api/pkg/api/response"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

var unauthorized = errors.New(errors.Input{
	Name:        "UnauthorizedWithLogoutError",
	StatusCode:  http.StatusUnauthorized,
	SendToSlack: errors.Bool(false),
	Message:     "Missing token in request.",
	Code:        "INVALID_JWT_TOKEN",
})

func Authenticated(c *gin.Context) {
	bearerToken := apicontext.BearerToken(c)
	if bearerToken == "" {
		apiresponse.Error(c, unauthorized)
		return
	}

	if len(strings.Split(bearerToken, ".")) != 3 {
		unauthorized.Message = "The token is badly formatted."
		apiresponse.Error(c, unauthorized)
		return
	}

	tokenClient := apicontext.TokenClient(c)
	unauthorized.Message = "Your access token is not valid, please login again."
	decoded, err := tokenClient.Decode(bearerToken)

	if err != nil {
		unauthorized.OriginalError = err.Error()
		apiresponse.Error(c, unauthorized)
		return
	}

	if _, ok := decoded.Meta["type"]; !ok {
		apiresponse.Error(c, unauthorized)
		return
	}

	c.Set(token.CtxDecodedKey, decoded)
	c.Next()
}
