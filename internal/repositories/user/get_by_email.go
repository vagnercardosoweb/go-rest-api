package user

import (
	"database/sql"

	"github.com/google/uuid"
)

type GetByEmailOutput struct {
	Id                uuid.UUID
	LoginBlockedUntil sql.NullTime `db:"login_blocked_until"`
	PasswordHash      string       `db:"password_hash"`
	Email             string
}

func (r *Repository) GetByEmail(email string) (*GetByEmailOutput, error) {
	var output GetByEmailOutput
	err := r.pgClient.QueryOne(&output, "SELECT id, email, password_hash, login_blocked_until FROM users WHERE LOWER(email) = LOWER($1) LIMIT 1", email)
	if err != nil {
		return nil, err
	}
	return &output, nil
}
