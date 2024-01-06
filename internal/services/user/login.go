package user

import (
	"github.com/vagnercardosoweb/go-rest-api/internal/repositories/user"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password_hash"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
	"net/http"
	"time"
)

type LoginSvc struct {
	tokenClient    token.Client
	userRepository user.RepositoryInterface
	passwordHash   password_hash.PasswordHash
}

func NewLoginSvc(userRepository user.RepositoryInterface, passwordHash password_hash.PasswordHash, tokenClient token.Client) *LoginSvc {
	return &LoginSvc{userRepository: userRepository, passwordHash: passwordHash, tokenClient: tokenClient}
}

func (s *LoginSvc) Execute(email, password string) (*token.Output, error) {
	userByEmail, err := s.userRepository.GetByEmail(email)
	if err != nil {
		return nil, s.invalidCredentialsError(err)
	}

	err = s.passwordHash.Compare(userByEmail.PasswordHash, password)
	if err != nil {
		return nil, s.invalidCredentialsError(err)
	}

	if userByEmail.LoginBlockedUntil.Valid && userByEmail.LoginBlockedUntil.Time.After(time.Now()) {
		return nil, errors.New(errors.Input{
			Message:    "Your access is blocked until: %s.",
			Arguments:  []any{userByEmail.LoginBlockedUntil.Time.Format("02/01/2006 at 15:04")},
			StatusCode: http.StatusUnauthorized,
		})
	}

	tokenOutput, err := s.tokenClient.Encode(&token.Input{Subject: userByEmail.Id.String()})
	if err != nil {
		return nil, err
	}

	return tokenOutput, nil
}

func (s *LoginSvc) invalidCredentialsError(originalError error) error {
	return errors.New(errors.Input{
		Message:       "Email/Password invalid",
		StatusCode:    http.StatusUnauthorized,
		OriginalError: originalError,
	})
}
