package user

import (
	"net/http"

	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

func (svc *Service) Login(email, password string) (string, error) {
	var resultToken string
	if email == "" || password == "" {
		return resultToken, errors.New(errors.Input{Message: "Email/Password is required", StatusCode: http.StatusBadRequest})
	}

	user, err := svc.userRepository.GetUserByEmailToLogin(email)
	if err != nil {
		return resultToken, errors.New(errors.Input{
			StatusCode: http.StatusUnauthorized,
			Message:    "Email/Password is invalid",
			Metadata:   errors.Metadata{"error": err.Error()},
		})
	}

	err = svc.passwordHash.Compare(user.PasswordHash, password)
	if err != nil {
		return resultToken, errors.New(errors.Input{
			Message:    "Email/Password is invalid",
			StatusCode: http.StatusUnauthorized,
			Metadata:   errors.Metadata{"error": err.Error()},
		})
	}

	if user.LoginBlockedUntil.Valid {
		return resultToken, errors.New(errors.Input{
			Message:   "Your login was blocked until: %s",
			Arguments: []any{user.LoginBlockedUntil.Time.UTC()},
		})
	}

	resultToken, err = svc.token.Encode(token.Input{Subject: user.ID.String()})
	if err != nil {
		return resultToken, errors.New(errors.Input{
			StatusCode: http.StatusInternalServerError,
			Metadata:   errors.Metadata{"error": err.Error()},
		})
	}

	return resultToken, nil
}
