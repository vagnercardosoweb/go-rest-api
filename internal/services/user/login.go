package user

import (
	"net/http"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/internal/events"
	"github.com/vagnercardosoweb/go-rest-api/internal/repositories/user"
	"github.com/vagnercardosoweb/go-rest-api/internal/types"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

type LoginSvc struct {
	tokenClient    token.Client
	passwordHash   password.PasswordHasher
	userRepository user.Repository
	eventManager   *events.EventManager
}

func NewLoginSvc(
	tokenClient token.Client,
	passwordHash password.PasswordHasher,
	userRepository user.Repository,
	eventManager *events.EventManager,
) *LoginSvc {
	return &LoginSvc{
		tokenClient:    tokenClient,
		passwordHash:   passwordHash,
		userRepository: userRepository,
		eventManager:   eventManager,
	}
}

func (s *LoginSvc) Execute(input *types.UserLoginInput) (*token.Output, error) {
	user, err := s.userRepository.GetByEmail(input.Email)
	if err != nil {
		return nil, s.invalidCredentialsError(err)
	}

	if err = s.passwordHash.Compare(user.PasswordHash, input.Password); err != nil {
		return nil, s.invalidCredentialsError(err)
	}

	if user.LoginBlockedUntil.Time.After(time.Now()) {
		return nil, errors.New(errors.Input{
			Message:    `Your access is blocked until "%s". Try again later.`,
			Arguments:  []any{user.LoginBlockedUntil.Time.Format("02/01/2006 at 15:04")},
			StatusCode: http.StatusUnauthorized,
		})
	}

	inputToken := &token.Input{Subject: user.Id.String()}
	outputToken, err := s.tokenClient.Encode(inputToken)

	if err != nil {
		return nil, err
	}

	s.eventManager.AfterLogin(events.AfterLoginInput{
		UserId: user.Id.String(),
	})

	return outputToken, nil
}

func (s *LoginSvc) invalidCredentialsError(originalError error) error {
	return errors.New(errors.Input{
		Message:       "Email/Password invalid",
		StatusCode:    http.StatusUnauthorized,
		OriginalError: originalError,
	})
}
