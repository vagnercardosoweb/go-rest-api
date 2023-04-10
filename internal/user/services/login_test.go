package user

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	userRepositories "github.com/vagnercardosoweb/go-rest-api/internal/user/repositories"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password_hash"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
	"github.com/vagnercardosoweb/go-rest-api/sqlc/store"
)

type UserServiceLoginTestSuite struct {
	suite.Suite
	passwordHashMock   *password_hash.PasswordHashMock
	userRepositoryMock *userRepositories.RepositoryMock
	tokenMock          *token.TokenMock
	svc                ServiceInterface
}

func (suite *UserServiceLoginTestSuite) SetupTest() {
	suite.passwordHashMock = password_hash.NewPasswordHashMock()
	suite.userRepositoryMock = userRepositories.NewRepositoryMock()
	suite.tokenMock = token.NewMock()
	suite.svc = NewService(suite.userRepositoryMock, suite.passwordHashMock, suite.tokenMock)
}

func (suite *UserServiceLoginTestSuite) TestLogin() {
	anyEmail := "user@mail.com"
	anyPassword := "any_password"

	suite.tokenMock.On("Encode", mock.Anything).Return("any_token", nil)
	suite.userRepositoryMock.On("GetUserByEmailToLogin", anyEmail).Return(&store.GetUserByEmailToLoginRow{PasswordHash: "hashed_password"}, nil)
	suite.passwordHashMock.On("Compare", "hashed_password", anyPassword).Return(nil)

	t := suite.T()
	loginToken, err := suite.svc.Login(anyEmail, anyPassword)

	assert.Nil(t, err)
	assert.Equal(t, loginToken, "any_token")

	suite.tokenMock.AssertExpectations(t)
	suite.tokenMock.AssertNumberOfCalls(t, "Encode", 1)

	suite.userRepositoryMock.AssertExpectations(t)
	suite.userRepositoryMock.AssertNumberOfCalls(t, "GetUserByEmailToLogin", 1)

	suite.passwordHashMock.AssertExpectations(t)
	suite.passwordHashMock.AssertNumberOfCalls(t, "Compare", 1)
}

func (suite *UserServiceLoginTestSuite) TestEmailAndPasswordIsEmpty() {
	loginToken, err := suite.svc.Login("", "")
	assert.NotNil(suite.T(), err)
	assert.Empty(suite.T(), loginToken)
}

func (suite *UserServiceLoginTestSuite) TestUserNotExist() {
	email := "user@mail.com"
	suite.userRepositoryMock.On("GetUserByEmailToLogin", email).Return(&store.GetUserByEmailToLoginRow{}, errors.New("any"))
	loginToken, err := suite.svc.Login(email, "any")
	assert.NotNil(suite.T(), err)
	assert.Empty(suite.T(), loginToken)
}

func (suite *UserServiceLoginTestSuite) TestComparePasswordError() {
	email := "user@mail.com"
	suite.userRepositoryMock.On("GetUserByEmailToLogin", email).Return(&store.GetUserByEmailToLoginRow{PasswordHash: "hashed_password"}, nil)
	suite.passwordHashMock.On("Compare", "hashed_password", "wrong_password").Return(errors.New("any"))
	_, err := suite.svc.Login(email, "wrong_password")
	assert.NotNil(suite.T(), err)
}

func (suite *UserServiceLoginTestSuite) TestLoginBlocked() {
	email := "user@mail.com"
	suite.passwordHashMock.On("Compare", mock.Anything, mock.Anything).Return(nil)
	suite.userRepositoryMock.On("GetUserByEmailToLogin", email).Return(&store.GetUserByEmailToLoginRow{LoginBlockedUntil: sql.NullTime{Time: time.Now().Add(time.Hour * 2), Valid: true}}, nil)
	_, err := suite.svc.Login(email, "any")
	assert.NotNil(suite.T(), err)
}

func (suite *UserServiceLoginTestSuite) TestErrorEncodeToken() {
	email := "user@mail.com"
	suite.passwordHashMock.On("Compare", mock.Anything, mock.Anything).Return(nil)
	suite.userRepositoryMock.On("GetUserByEmailToLogin", email).Return(&store.GetUserByEmailToLoginRow{}, nil)
	suite.tokenMock.On("Encode", mock.Anything).Return("", errors.New("any"))
	token, err := suite.svc.Login(email, "any")
	assert.NotNil(suite.T(), err)
	assert.Empty(suite.T(), token)
}

func TestUserServiceLoginTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceLoginTestSuite))
}
