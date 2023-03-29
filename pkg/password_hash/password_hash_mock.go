package password_hash

import "github.com/stretchr/testify/mock"

type PasswordHashMock struct {
	mock.Mock
}

func NewPasswordHashMock() *PasswordHashMock {
	return &PasswordHashMock{}
}

func (m *PasswordHashMock) Compare(hashedPassword string, plainPassword string) error {
	args := m.Called(hashedPassword, plainPassword)
	return args.Error(0)
}

func (m *PasswordHashMock) Create(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}
