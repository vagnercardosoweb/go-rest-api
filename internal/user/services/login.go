package user

import (
	"net/http"

	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

func makeUnauthorizedError(err error) error {
	return errors.New(errors.Input{
		StatusCode: http.StatusUnauthorized,
		Message:    "Email/Password is invalid",
		Metadata:   errors.Metadata{"error": err.Error()},
	})
}

func (svc *Service) Login(email, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New(errors.Input{
			Message:    "Email/Password is required",
			StatusCode: http.StatusBadRequest,
		})
	}

	user, err := svc.userRepository.GetUserByEmailToLogin(email)
	if err != nil {
		return "", makeUnauthorizedError(err)
	}

	err = svc.passwordHash.Compare(user.PasswordHash, password)
	if err != nil {
		return "", makeUnauthorizedError(err)
	}

	if user.LoginBlockedUntil.Valid {
		return "", errors.New(errors.Input{
			Message:    "Your login was blocked until: %s",
			Arguments:  []any{user.LoginBlockedUntil.Time.UTC()},
			StatusCode: http.StatusForbidden,
		})
	}

	resultToken, err := svc.token.Encode(token.Input{Subject: user.ID.String()})
	if err != nil {
		return "", errors.New(errors.Input{
			StatusCode: http.StatusInternalServerError,
			Metadata:   errors.Metadata{"error": err.Error()},
		})
	}

	return resultToken, nil
}
