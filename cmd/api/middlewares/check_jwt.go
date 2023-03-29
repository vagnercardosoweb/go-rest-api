package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

func CheckJwt(c *gin.Context) {
	headerToken := c.GetString(config.AuthHeaderToken)

	if headerToken == env.Required("TOKEN_DEFAULT") {
		c.Next()
		return
	}

	unauthorized := errors.New(errors.Input{
		Message:    "Missing token in request.",
		StatusCode: http.StatusUnauthorized,
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

	payload, err := token.NewJwt(nil).Decode(headerToken)
	if err != nil {
		unauthorized.Message = "Your access token is not valid, please login again."
		unauthorized.AddMetadata("error", err.Error())
		c.AbortWithError(unauthorized.StatusCode, unauthorized)
		return
	}

	c.Set(config.TokenPayloadCtxKey, payload)
	c.Next()
}
