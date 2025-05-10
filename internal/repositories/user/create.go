package user

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

type CreateInput struct {
	Name              string
	Email             string
	PasswordHash      string
	CodeToInvite      string
	ConfirmedEmailAt  sql.NullTime
	LoginBlockedUntil sql.NullTime
	Birthdate         time.Time
}

type CreateOutput struct {
	CreateInput
	Id uuid.UUID
}

const createQuery = `
	INSERT INTO
		users (
			id,
			NAME,
			email,
			password_hash,
			code_to_invite,
			confirmed_email_at,
			login_blocked_until,
			birth_date
		)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8);
`

func (r *Repository) Create(input *CreateInput) (*CreateOutput, error) {
	id := uuid.New()

	_, err := r.pgClient.Exec(
		createQuery,
		id,
		input.Name,
		input.Email,
		input.PasswordHash,
		input.CodeToInvite,
		input.ConfirmedEmailAt,
		input.LoginBlockedUntil,
		input.Birthdate,
	)

	if err != nil {
		return nil, errors.FromSql(err)
	}

	return &CreateOutput{
		CreateInput: *input,
		Id:          id,
	}, nil
}
