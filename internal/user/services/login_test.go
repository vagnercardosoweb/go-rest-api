package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	userRepositories "github.com/vagnercardosoweb/go-rest-api/internal/user/repositories"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password_hash"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
	"github.com/vagnercardosoweb/go-rest-api/sqlc/store"
)

func TestUserServices_Login(t *testing.T) {
	passwordHashMock := password_hash.NewPasswordHashMock()
	userRepositoryMock := userRepositories.NewRepositoryMock()
	tokenMock := token.NewMock()

	anyEmail := "any_email@mail.com"
	anyPassword := "any_password"

	tokenMock.On("Encode", mock.Anything).Return("any_token", nil)
	userRepositoryMock.On("GetUserByEmailToLogin", anyEmail).Return(&store.GetUserByEmailToLoginRow{PasswordHash: "hashed_password"}, nil)
	passwordHashMock.On("Compare", "hashed_password", anyPassword).Return(nil)

	userService := NewService(userRepositoryMock, passwordHashMock, tokenMock)
	loginToken, err := userService.Login(anyEmail, anyPassword)

	assert.Nil(t, err)
	assert.Equal(t, loginToken, "any_token")

	tokenMock.AssertExpectations(t)
	tokenMock.AssertNumberOfCalls(t, "Encode", 1)

	userRepositoryMock.AssertExpectations(t)
	userRepositoryMock.AssertNumberOfCalls(t, "GetUserByEmailToLogin", 1)

	passwordHashMock.AssertExpectations(t)
	passwordHashMock.AssertNumberOfCalls(t, "Compare", 1)

}
