package password

import (
	"context"
	"fmt"
)

type PasswordHasher interface {
	Compare(hashedPassword string, plainPassword string) error
	Create(password string) (string, error)
}

const CtxKey = "PasswordHasherKey"

func FromCtx(c context.Context) PasswordHasher {
	value, exists := c.Value(CtxKey).(PasswordHasher)

	if !exists {
		panic(fmt.Errorf(`context key "%s" does not exist`, CtxKey))
	}

	return value
}
