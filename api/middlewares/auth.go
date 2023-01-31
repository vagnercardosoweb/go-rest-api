package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	jwt "github.com/vagnercardosoweb/go-rest-api/pkg/jwt"
)

func Auth(c *gin.Context) {
	token := c.GetString(config.AuthTokenContextKey)
	unauthorizedError := errors.NewUnauthorized(errors.Input{Message: "Missing token in request."})

	if token == "" {
		c.AbortWithError(unauthorizedError.StatusCode, unauthorizedError)
		return
	}

	if token == env.Get("JWT_SECRET_KEY") {
		c.Next()
		return
	}

	if len(strings.Split(token, ".")) != 3 {
		unauthorizedError.Message = "The token is badly formatted."
		c.AbortWithError(unauthorizedError.StatusCode, unauthorizedError)
		return
	}

	payload, err := jwt.Verify(token)

	if err != nil {
		unauthorizedError.Message = err.Error()
		c.AbortWithError(unauthorizedError.StatusCode, unauthorizedError)
		return
	}

	// userId := payload.Sub
	// dbConnection := c.MustGet(config.DbConnectionContextKey).(*postgres.Connection)
	// var user map[string]any
	// dbConnection.Query(user, "select * from users where id = $1", userId)

	c.Set(config.JwtPayloadContextKey, payload)
	c.Next()
}
