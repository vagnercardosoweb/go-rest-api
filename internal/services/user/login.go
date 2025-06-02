package user

import (
	"net/http"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/internal/events"
	"github.com/vagnercardosoweb/go-rest-api/internal/repositories/user"
	"github.com/vagnercardosoweb/go-rest-api/internal/types"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

type LoginSvc struct {
	pgClient       *postgres.Client
	eventManager   *events.Manager
	userRepository user.Repository
	passwordHash   password.PasswordHasher
	tokenClient    token.Client
}

func NewLoginSvc(
	pgClient *postgres.Client,
	tokenClient token.Client,
	passwordHash password.PasswordHasher,
	eventManager *events.Manager,
) *LoginSvc {
	return &LoginSvc{
		pgClient:       pgClient,
		tokenClient:    tokenClient,
		userRepository: user.New(pgClient),
		passwordHash:   passwordHash,
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
			StatusCode: http.StatusUnauthorized,
			Message:    `Your access is blocked until "%s". Try again later.`,
			Arguments:  []any{user.LoginBlockedUntil.Time.Format("02/01/2006 at 15:04")},
		})
	}

	inputToken := &token.Input{Subject: user.Id.String()}
	outputToken, err := s.tokenClient.Encode(inputToken)

	if err != nil {
		return nil, err
	}

	s.eventManager.SendOnUserLogin(events.OnUserLoginInput{
		UserId:    user.Id.String(),
		RequestId: s.pgClient.Logger().GetId(),
		UserAgent: input.UserAgent,
		IpAddress: input.IpAddress,
	})

	return outputToken, nil
}

func (s *LoginSvc) invalidCredentialsError(originalError error) error {
	return errors.New(errors.Input{
		StatusCode:    http.StatusUnauthorized,
		Message:       "user.invalidCredentials",
		OriginalError: originalError,
	})
}
