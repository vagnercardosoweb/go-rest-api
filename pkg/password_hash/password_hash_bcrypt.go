package password_hash

import (
	"golang.org/x/crypto/bcrypt"
)

type PasswordHashBcrypt struct {
	cost int
}

func NewPasswordHashBcrypt() *PasswordHashBcrypt {
	return &PasswordHashBcrypt{12}
}

func (b *PasswordHashBcrypt) WithCost(cost int) *PasswordHashBcrypt {
	b.cost = cost
	return b
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
