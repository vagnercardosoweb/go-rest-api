package user

import (
	"database/sql"

	"github.com/google/uuid"
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

func (r *Repository) GetByEmail(email string) (*GetByEmailOutput, error) {
	var output GetByEmailOutput

	err := r.pgClient.QueryOne(&output, getByEmailQuery, email)
	if err != nil {
		return nil, err
	}

	return &output, nil
}
