package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/jwt"
)

func Auth(c *gin.Context) {
	token := c.GetString(config.AuthTokenCtxKey)

	if token == env.Required("JWT_SECRET_KEY") {
		c.Next()
		return
	}

	unauthorized := errors.NewUnauthorized(errors.Input{Message: "Missing token in request."})

	if token == "" {
		c.AbortWithError(unauthorized.StatusCode, unauthorized)
		return
	}

	if len(strings.Split(token, ".")) != 3 {
		unauthorized.Message = "The token is badly formatted."
		c.AbortWithError(unauthorized.StatusCode, unauthorized)
		return
	}

	payload, err := jwt.Verify(token)

	if err != nil {
		unauthorized.Message = "Your access token is not valid, please login again."
		unauthorized.AddMetadata("error", err.Error())
		c.AbortWithError(unauthorized.StatusCode, unauthorized)
		return
	}

	c.Set(config.JwtPayloadCtxKey, payload)
	c.Next()
}
