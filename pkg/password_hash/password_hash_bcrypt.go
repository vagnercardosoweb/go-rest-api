package password_hash

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHashBcrypt struct {
	cost int
}

func NewPasswordHashBcrypt(cost int) PasswordHash {
	return &PasswordHashBcrypt{cost}
}

func (b *PasswordHashBcrypt) Create(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err == nil {
		return string(bytes), nil
	}
	return "", err
}

func (b *PasswordHashBcrypt) Compare(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
