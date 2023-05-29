package token

import (
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func NewMock() *Mock {
	return &Mock{}
}

func (m *Mock) Encode(input Input) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

func (m *Mock) Decode(token string) (*Output, error) {
	args := m.Called(token)
	return args.Get(0).(*Output), args.Error(1)
}
