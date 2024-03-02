package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/vagnercardosoweb/go-rest-api/internal/events"
	repository "github.com/vagnercardosoweb/go-rest-api/internal/repositories/user"
	service "github.com/vagnercardosoweb/go-rest-api/internal/services/user"
	"github.com/vagnercardosoweb/go-rest-api/internal/types"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password_hash"
)

func Login(c *gin.Context) any {
	var input types.UserLoginInput
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

	eventManager := events.GetFromGinCtx(c)
	eventManager.Dispatch(events.MakeAfterLogin(result.Subject))

	return types.UserLoginOutput{
		AccessToken: result.Token,
		ExpiresIn:   result.ExpiresAt,
		TokenType:   "Bearer",
	}
}
