package auth_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	repositories "github.com/vagnercardosoweb/go-rest-api/internal/auth/repositories"
	services "github.com/vagnercardosoweb/go-rest-api/internal/auth/services"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password_hash"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

type Input struct {
	Email    string
	Password string
}

func Login(c *gin.Context) any {
	var input Input
	if err := c.ShouldBindBodyWith(&input, binding.JSON); err != nil {
		return err
	}

	svc := services.New(
		repositories.NewPostgres(
			c.MustGet(config.PgConnectCtxKey).(*postgres.Connection),
			c.Request.Context(),
		),
		password_hash.NewBcrypt(),
		token.NewJwt(),
	)

	result, err := svc.Login(input.Email, input.Password)
	if err != nil {
		return err
	}

	return map[string]string{
		"token": result,
	}
}
