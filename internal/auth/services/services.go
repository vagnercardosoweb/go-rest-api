package auth_services

import (
	"net/http"
	"time"

	authRepositories "github.com/vagnercardosoweb/go-rest-api/internal/auth/repositories"

	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password_hash"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

type Interface interface {
	Login(email, password string) (string, error)
}

type service struct {
	token          token.Token
	authRepository authRepositories.Interface
	passwordHash   password_hash.PasswordHash
}

func New(authRepository authRepositories.Interface, passwordHash password_hash.PasswordHash, token token.Token) Interface {
	return &service{authRepository: authRepository, passwordHash: passwordHash, token: token}
}

func (svc *service) Login(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New(errors.Input{Message: "Email/Password is required"})
	}

	unauthorizedError := errors.New(errors.Input{
		Message:    "Email/Password invalid",
		StatusCode: http.StatusUnauthorized,
	})

	user, err := svc.authRepository.GetUserByEmail(email)
	if err != nil {
		unauthorizedError.OriginalError = err.Error()
		return "", unauthorizedError
	}

	err = svc.passwordHash.Compare(user.PasswordHash, password)
	if err != nil {
		unauthorizedError.OriginalError = err.Error()
		return "", unauthorizedError
	}

	if user.LoginBlockedUntil.Valid && user.LoginBlockedUntil.Time.After(time.Now()) {
		return "", errors.New(errors.Input{
			Message:   "Your access is blocked until: %s.",
			Arguments: []any{user.LoginBlockedUntil.Time.Format("02/01/2006 at 15:04")},
		})
	}

	jwt, err := svc.token.Encode(token.Input{Subject: user.ID.String()})
	if err != nil {
		return "", err
	}

	return jwt, nil
}
