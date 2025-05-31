package user

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/vagnercardosoweb/go-rest-api/internal/events"
	userRepository "github.com/vagnercardosoweb/go-rest-api/internal/repositories/user"
	userService "github.com/vagnercardosoweb/go-rest-api/internal/services/user"
	"github.com/vagnercardosoweb/go-rest-api/internal/types"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func Login(c *gin.Context) any {
	var input *types.UserLoginInput

	if err := c.ShouldBindBodyWith(&input, binding.JSON); err != nil {
		return errors.FromTranslator(err, apicontext.ValidatorTranslator(c))
	}

	loginSvc := userService.NewLoginSvc(
		apicontext.TokenClient(c),
		apicontext.PasswordHasher(c),
		userRepository.New(apicontext.PgClient(c)),
		events.FromGin(c),
	)

	result, err := loginSvc.Execute(input)
	if err != nil {
		return err
	}

	return types.UserLoginOutput{
		AccessToken: result.Token,
		ExpiresIn:   result.ExpiresAt,
		TokenType:   "Bearer",
	}
}
