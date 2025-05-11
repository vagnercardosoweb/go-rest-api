package password

import (
	"net/http"

	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
	cost int
}

func NewBcrypt() *Bcrypt {
	return &Bcrypt{12}
}

func (b *Bcrypt) WithCost(cost int) *Bcrypt {
	b.cost = cost
	return b
}

func (b *Bcrypt) Create(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err == nil {
		return string(bytes), nil
	}
	return "", err
}

func (b *Bcrypt) Compare(hashedPassword string, password string) error {
	if len(hashedPassword) == 0 {
		return errors.New(errors.Input{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Hashed password is empty",
		})
	}

	if len(password) == 0 {
		return errors.New(errors.Input{
			StatusCode: http.StatusUnprocessableEntity,
			Message:    "Password is empty",
		})
	}

	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
}
