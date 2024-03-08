package user_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/vagnercardosoweb/go-rest-api/internal/repositories/user"
	"github.com/vagnercardosoweb/go-rest-api/internal/types"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password_hash"
	"github.com/vagnercardosoweb/go-rest-api/tests"
)

type LoginSuite struct {
	tests.RestApiSuite
	passwordHash   password_hash.PasswordHash
	userRepository user.RepositoryInterface
	input          *types.UserLoginInput
}

func (suite *LoginSuite) SetupSuite() {
	suite.RestApiSuite.SetupSuite()

	suite.passwordHash = password_hash.NewBcrypt()
	suite.userRepository = user.NewRepository(suite.PgClient)

	suite.input = &types.UserLoginInput{
		Email:    "test@local.dev",
		Password: "12345678",
	}
}

func (suite *LoginSuite) SetupTest() {
	passwordHash, err := suite.passwordHash.Create(suite.input.Password)
	suite.Require().Nil(err)

	_, err = suite.userRepository.Create(&user.CreateInput{
		Name:         "Test User",
		Email:        suite.input.Email,
		PasswordHash: passwordHash,
		Birthdate:    time.Date(1994, time.December, 15, 0, 0, 0, 0, time.UTC),
		CodeToInvite: "ANY_CODE",
	})

	suite.Require().Nil(err)
}

func (suite *LoginSuite) TearDownTest() {
	_ = suite.PgClient.TruncateTable("users")
}

func (suite *LoginSuite) makeRecorder(input any) *httptest.ResponseRecorder {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(input)
	rr := suite.RestApi.TestRequest(httptest.NewRequest(http.MethodPost, "/login", body))
	return rr
}

func (suite *LoginSuite) TestSuccess() {
	rr := suite.makeRecorder(suite.input)
	suite.Require().Equal(http.StatusOK, rr.Code)

	var output struct {
		Data types.UserLoginOutput `json:"data"`
	}

	_ = json.NewDecoder(rr.Body).Decode(&output)

	suite.Require().NotEmpty(output.Data.AccessToken)
	suite.Require().True(len(strings.Split(output.Data.AccessToken, ".")) == 3)
	suite.Require().Equal(output.Data.TokenType, "Bearer")
	suite.Require().NotEmpty(output.Data.ExpiresIn)

}

func (suite *LoginSuite) TestNotFound() {
	input := suite.input
	input.Email = "not_found@local.dev"

	rr := suite.makeRecorder(input)

	var e errors.Input
	_ = json.NewDecoder(rr.Body).Decode(&e)

	suite.Require().Equal(http.StatusUnauthorized, rr.Code)
	suite.Require().Equal(e.Message, "Email/Password invalid")
}

func (suite *LoginSuite) TestInvalidPassword() {
	input := suite.input
	input.Password = "invalid_password"

	rr := suite.makeRecorder(input)

	var e errors.Input
	_ = json.NewDecoder(rr.Body).Decode(&e)

	suite.Require().Equal(http.StatusUnauthorized, rr.Code)
	suite.Require().Equal(e.Message, "Email/Password invalid")
}

func (suite *LoginSuite) TestBlockedUntil() {
	_, _ = suite.PgClient.Exec("UPDATE users SET login_blocked_until = NOW() + INTERVAL '1 hour' WHERE email = $1", suite.input.Email)

	rr := suite.makeRecorder(suite.input)

	var e errors.Input
	_ = json.NewDecoder(rr.Body).Decode(&e)

	suite.Require().Equal(http.StatusUnauthorized, rr.Code)
	suite.Require().Equal(
		e.Message,
		fmt.Sprintf(
			`Your access is blocked until "%s". Try again later.`,
			time.Now().Add(time.Hour).Format("02/01/2006 at 15:04"),
		),
	)
}

func TestLoginSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	suite.Run(t, new(LoginSuite))
}
