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
	"github.com/vagnercardosoweb/go-rest-api/pkg/password"
	"github.com/vagnercardosoweb/go-rest-api/tests"
)

type LoginTest struct {
	tests.RestApiSuite
	passwordHash   password.PasswordHasher
	userRepository user.Repository
	validInput     types.UserLoginInput
}

func (t *LoginTest) createRecorder(input any) *httptest.ResponseRecorder {
	body := new(bytes.Buffer)
	_ = json.NewEncoder(body).Encode(input)

	request := httptest.NewRequest(http.MethodPost, "/login", body)
	rr := t.RestApi.TestRequest(request)

	return rr
}

func (t *LoginTest) SetupSuite() {
	t.RestApiSuite.SetupSuite()

	t.passwordHash = password.NewBcrypt()
	t.userRepository = user.New(t.PgClient)

	t.validInput = types.UserLoginInput{
		Email:    "test@test.local",
		Password: "12345678",
	}
}

func (t *LoginTest) SetupTest() {
	passwordHash, err := t.passwordHash.Create(t.validInput.Password)
	t.Require().Nil(err)

	_, err = t.userRepository.Create(&user.CreateInput{
		Name:         "Test User",
		Email:        t.validInput.Email,
		PasswordHash: passwordHash,
		Birthdate:    time.Date(1994, time.December, 15, 0, 0, 0, 0, time.UTC),
		CodeToInvite: "ANY_CODE",
	})

	t.Require().Nil(err)
}

func (t *LoginTest) TearDownTest() {
	_ = t.PgClient.TruncateTable("users")
}

func (t *LoginTest) checkLastLogin() {
	lastLogin := new(struct {
		LastLoginAt    *time.Time `db:"last_login_at"`
		LastLoginAgent *string    `db:"last_login_agent"`
		LastLoginIp    *string    `db:"last_login_ip"`
	})

	query := `SELECT "last_login_at", "last_login_agent", "last_login_ip" FROM "users" WHERE LOWER("email") = LOWER($1);`
	err := t.PgClient.QueryRow(lastLogin, query, t.validInput.Email)
	t.Require().Nil(err)

	t.Require().NotNil(lastLogin.LastLoginAt, "Last login time should not be nil")
	t.Require().NotNil(lastLogin.LastLoginAgent, "Last login agent should not be nil")
	t.Require().NotNil(lastLogin.LastLoginIp, "Last login IP should not be nil")
}

func (t *LoginTest) TestSuccess() {
	rr := t.createRecorder(t.validInput)
	t.Require().Equal(http.StatusOK, rr.Code)

	var output types.UserLoginOutput
	_ = json.NewDecoder(rr.Body).Decode(&output)

	t.Require().NotEmpty(output.AccessToken)
	t.Require().True(len(strings.Split(output.AccessToken, ".")) == 3)
	t.Require().Equal(output.TokenType, "Bearer")
	t.Require().NotEmpty(output.ExpiresIn)

	t.checkLastLogin()
}

func (t *LoginTest) TestNotFound() {
	rr := t.createRecorder(types.UserLoginInput{
		Email:    "not_found@test.local",
		Password: t.validInput.Password,
	})

	var e errors.Input
	_ = json.NewDecoder(rr.Body).Decode(&e)

	t.Require().Equal(http.StatusUnauthorized, rr.Code)
	t.Require().Equal(e.Message, "user.invalidCredentials")
}

func (t *LoginTest) TestInvalidPassword() {
	rr := t.createRecorder(types.UserLoginInput{
		Email:    t.validInput.Email,
		Password: "invalid_password",
	})

	var e errors.Input
	_ = json.NewDecoder(rr.Body).Decode(&e)

	t.Require().Equal(http.StatusUnauthorized, rr.Code)
	t.Require().Equal(e.Message, "user.invalidCredentials")
}

func (t *LoginTest) TestBlockedUntil() {
	_, _ = t.PgClient.Exec(
		`UPDATE "users" SET "login_blocked_until" = NOW() + INTERVAL '1 HOUR' WHERE "email" = $1;`,
		t.validInput.Email,
	)

	rr := t.createRecorder(t.validInput)

	var e errors.Input
	_ = json.NewDecoder(rr.Body).Decode(&e)

	t.Require().Equal(http.StatusUnauthorized, rr.Code)
	t.Require().Equal(
		e.Message,
		fmt.Sprintf(
			`Your access is blocked until "%s". Try again later.`,
			time.Now().Add(time.Hour).Format("02/01/2006 at 15:04"),
		),
	)
}

func TestUserLogin(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	suite.Run(t, new(LoginTest))
}
