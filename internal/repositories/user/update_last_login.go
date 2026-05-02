package user

import (
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

type UpdateLastLoginInput struct {
	UserId    string
	IpAddress string
	UserAgent string
}

const updateLastLoginQuery = `UPDATE "users"
SET
	"last_login_at" = NOW(),
	"last_login_agent" = $2,
	"last_login_ip" = $1
WHERE
	"id" = $3;`

func (r *instance) UpdateLastLogin(input *UpdateLastLoginInput) error {
	_, err := r.pgClient.Exec(
		updateLastLoginQuery,
		input.IpAddress,
		input.UserAgent,
		input.UserId,
	)

	if err != nil {
		return errors.FromSql(err)
	}

	return nil
}
