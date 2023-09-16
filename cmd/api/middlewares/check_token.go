package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func CheckToken(c *gin.Context) {
	headerToken := c.GetString(config.AuthHeaderTokenCtx)
	unauthorized := errors.New(errors.Input{
		StatusCode:  http.StatusUnauthorized,
		SendToSlack: errors.Bool(false),
		Message:     "Missing token in request.",
		Code:        "INVALID_JWT_TOKEN",
	})

	if headerToken == "" {
		c.AbortWithError(unauthorized.StatusCode, unauthorized)
		return
	}

	if len(strings.Split(headerToken, ".")) != 3 {
		unauthorized.Message = "The token is badly formatted."
		c.AbortWithError(unauthorized.StatusCode, unauthorized)
		return
	}

	token := config.GetTokenFromCtx(c)
	decoded, err := token.Decode(headerToken)

	if err != nil {
		unauthorized.Message = "Your access token is not valid, please login again."
		unauthorized.OriginalError = err.Error()
		c.AbortWithError(unauthorized.StatusCode, unauthorized)
		return
	}

	if _, ok := decoded.Meta["type"]; !ok {
		unauthorized.Message = "Your access token is not valid, please login again."
		c.AbortWithError(unauthorized.StatusCode, unauthorized)
		return
	}

	c.Set(config.TokenDecodedCtxKey, decoded)
	c.Next()
}