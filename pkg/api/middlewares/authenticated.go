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

func buildUnauthorizedError(message string) error {
	return errors.New(errors.Input{
		Name:       "UnauthorizedWithLogoutError",
		StatusCode: http.StatusUnauthorized,
		Message:    message,
		Code:       "INVALID_JWT_TOKEN",
	})
}

func Authenticated(c *gin.Context) {
	bearerToken := apicontext.BearerToken(c)
	if bearerToken == "" {
		unauthorizedError := buildUnauthorizedError("Missing token in request.")
		apiresponse.Error(c, unauthorizedError)
		return
	}

	if len(strings.Split(bearerToken, ".")) != 3 {
		unauthorizedError := buildUnauthorizedError("The token is badly formatted.")
		apiresponse.Error(c, unauthorizedError)
		return
	}

	tokenClient := apicontext.TokenClient(c)
	decoded, err := tokenClient.Decode(bearerToken)

	if err != nil {
		unauthorizedError := buildUnauthorizedError("Your access token is not valid, please login again.")
		apiresponse.Error(c, unauthorizedError)
		return
	}

	c.Set(token.CtxDecodedKey, decoded)
	c.Next()
}
