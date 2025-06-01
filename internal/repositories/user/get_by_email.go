package user

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

type GetByEmailOutput struct {
	Id                uuid.UUID
	PasswordHash      string       `db:"password_hash"`
	LoginBlockedUntil sql.NullTime `db:"login_blocked_until"`
	Email             string
}

const getByEmailQuery = `
	SELECT
		"id",
		"email",
		"password_hash",
		"login_blocked_until"
	FROM
		"users"
	WHERE
		LOWER("email") = LOWER($1)
	LIMIT
		1;
`

func (r *instance) GetByEmail(email string) (*GetByEmailOutput, error) {
	output := new(GetByEmailOutput)

	err := r.pgClient.QueryRow(output, getByEmailQuery, email)
	if err != nil {
		return nil, errors.FromSql(err, "user.notFoundByEmail", email)
	}

	return output, nil
}
