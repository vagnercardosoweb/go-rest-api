package token

import (
	"github.com/stretchr/testify/mock"
)

type TokenMock struct {
	mock.Mock
}

func NewMock() *TokenMock {
	return &TokenMock{}
}

func (m *TokenMock) Encode(input Input) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

func (m *TokenMock) Decode(token string) (*Output, error) {
	args := m.Called(token)
	return args.Get(0).(*Output), args.Error(1)
}
