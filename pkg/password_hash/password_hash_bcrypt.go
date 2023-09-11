package password_hash

import (
	"errors"

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

func (b *Bcrypt) Compare(hashedPassword string, plainPassword string) error {
	if len(hashedPassword) == 0 {
		return errors.New("hashed password is empty")
	}
	if len(plainPassword) == 0 {
		return errors.New("plain password is empty")
	}
	return bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(plainPassword),
	)
}
