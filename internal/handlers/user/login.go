package user

import (
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/internal/events"
	"github.com/vagnercardosoweb/go-rest-api/internal/services/user"
	"github.com/vagnercardosoweb/go-rest-api/internal/types"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func Login(c *gin.Context) any {
	input := new(types.UserLoginInput)

	if err := c.ShouldBindBodyWithJSON(input); err != nil {
		return errors.FromTranslator(err, apicontext.ValidatorTranslator(c))
	}

	loginSvc := user.NewLoginSvc(
		apicontext.PgClient(c),
		apicontext.TokenClient(c),
		apicontext.PasswordHasher(c),
		events.FromGin(c),
	)

	input.IpAddress = c.ClientIP()
	input.UserAgent = c.GetHeader("User-Agent")

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
