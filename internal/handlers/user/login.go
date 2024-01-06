package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	repository "github.com/vagnercardosoweb/go-rest-api/internal/repositories/user"
	service "github.com/vagnercardosoweb/go-rest-api/internal/services/user"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password_hash"
)

type Input struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

func Login(c *gin.Context) any {
	var input Input
	if err := c.ShouldBindBodyWith(&input, binding.JSON); err != nil {
		return errors.FromBindJson(err, utils.GetValidateTranslator(c))
	}

	loginSvc := service.NewLoginSvc(
		repository.NewRepository(utils.GetPgClient(c)),
		password_hash.NewBcrypt(),
		utils.GetTokenClient(c),
	)

	result, err := loginSvc.Execute(input.Email, input.Password)
	if err != nil {
		return err
	}

	return map[string]any{
		"accessToken": result.Token,
		"expiresIn":   result.ExpiresAt,
		"tokenType":   "Bearer",
	}
}
