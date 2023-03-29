package user

import (
	"github.com/stretchr/testify/mock"
	"github.com/vagnercardosoweb/go-rest-api/sqlc/store"
)

type RepositoryMock struct {
	mock.Mock
}

func NewRepositoryMock() *RepositoryMock {
	return &RepositoryMock{}
}

func (m *RepositoryMock) GetUserByEmailToLogin(email string) (*store.GetUserByEmailToLoginRow, error) {
	args := m.Called(email)
	return args.Get(0).(*store.GetUserByEmailToLoginRow), args.Error(1)
}
