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
	input          types.UserLoginInput
}

func (ls *LoginSuite) SetupSuite() {
	ls.RestApiSuite.SetupSuite()

	ls.passwordHash = password_hash.NewBcrypt()
	ls.userRepository = user.NewRepository(ls.PgClient)

	ls.input = types.UserLoginInput{
		Email:    "test@local.dev",
		Password: "12345678",
	}
}

func (ls *LoginSuite) SetupTest() {
	passwordHash, err := ls.passwordHash.Create(ls.input.Password)
	ls.Require().Nil(err)

	_, err = ls.userRepository.Create(&user.CreateInput{
		Name:         "Test User",
		Email:        ls.input.Email,
		PasswordHash: passwordHash,
		Birthdate:    time.Date(1994, time.December, 15, 0, 0, 0, 0, time.UTC),
		CodeToInvite: "ANY_CODE",
	})

	ls.Require().Nil(err)
}

func (ls *LoginSuite) TearDownTest() {
	_ = ls.PgClient.TruncateTable("users")
}

func (ls *LoginSuite) createResponseRecorder(input any) *httptest.ResponseRecorder {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(input)
	rr := ls.RestApi.TestRequest(httptest.NewRequest(http.MethodPost, "/login", body))
	return rr
}

func (ls *LoginSuite) TestSuccess() {
	rr := ls.createResponseRecorder(ls.input)
	ls.Require().Equal(http.StatusOK, rr.Code)

	var output struct {
		Data types.UserLoginOutput `json:"data"`
	}

	_ = json.NewDecoder(rr.Body).Decode(&output)

	ls.Require().NotEmpty(output.Data.AccessToken)
	ls.Require().True(len(strings.Split(output.Data.AccessToken, ".")) == 3)
	ls.Require().Equal(output.Data.TokenType, "Bearer")
	ls.Require().NotEmpty(output.Data.ExpiresIn)

}

func (ls *LoginSuite) TestNotFound() {
	ls.input.Email = "not_found@local.dev"
	rr := ls.createResponseRecorder(ls.input)

	var e errors.Input
	_ = json.NewDecoder(rr.Body).Decode(&e)

	ls.Require().Equal(http.StatusUnauthorized, rr.Code)
	ls.Require().Equal(e.Message, "Email/Password invalid")
}

func (ls *LoginSuite) TestInvalidPassword() {
	ls.input.Password = "invalid_password"
	rr := ls.createResponseRecorder(ls.input)

	var e errors.Input
	_ = json.NewDecoder(rr.Body).Decode(&e)

	ls.Require().Equal(http.StatusUnauthorized, rr.Code)
	ls.Require().Equal(e.Message, "Email/Password invalid")
}

func (ls *LoginSuite) TestBlockedUntil() {
	_, _ = ls.PgClient.Exec(
		"UPDATE users SET login_blocked_until = NOW() + INTERVAL '1 hour' WHERE email = $1",
		ls.input.Email,
	)

	rr := ls.createResponseRecorder(ls.input)

	var e errors.Input
	_ = json.NewDecoder(rr.Body).Decode(&e)

	ls.Require().Equal(http.StatusUnauthorized, rr.Code)
	ls.Require().Equal(
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
