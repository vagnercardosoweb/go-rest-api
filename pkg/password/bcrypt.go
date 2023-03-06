package password

import (
	"golang.org/x/crypto/bcrypt"
)

type Bcrypt struct {
	cost int
}

func NewBcrypt(cost int) Password {
	return &Bcrypt{cost}
}

func (i *Bcrypt) Create(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), i.cost)
	if err == nil {
		return string(bytes), nil
	}
	return "", err
}

func (i *Bcrypt) Compare(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
