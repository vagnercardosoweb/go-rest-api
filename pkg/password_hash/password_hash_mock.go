package password_hash

import "github.com/stretchr/testify/mock"

type Mock struct {
	mock.Mock
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Compare(hashedPassword string, plainPassword string) error {
	args := m.Called(hashedPassword, plainPassword)
	return args.Error(0)
}

func (m *Mock) Create(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}
