package auth_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	repositories "github.com/vagnercardosoweb/go-rest-api/internal/auth/repositories"
	services "github.com/vagnercardosoweb/go-rest-api/internal/auth/services"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password_hash"
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

	token := config.GetTokenFromCtx(c)
	pgClient := config.GetPgClientFromCtx(c)

	svc := services.New(
		repositories.NewPostgres(pgClient, c.Request.Context()),
		password_hash.NewBcrypt(),
		token,
	)

	result, err := svc.Login(input.Email, input.Password)
	if err != nil {
		return err
	}

	return map[string]string{
		"token": result,
	}
}
