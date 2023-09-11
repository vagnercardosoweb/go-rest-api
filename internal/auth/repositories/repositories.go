package auth_repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
)

type Interface interface {
	GetUserByEmail(email string) (*GetUserByEmailOutput, error)
}

type pg struct {
	ctx context.Context
	db  *postgres.Client
}

func NewPostgres(db *postgres.Client, ctx context.Context) Interface {
	return &pg{db: db, ctx: ctx}
}

type GetUserByEmailOutput struct {
	ID                uuid.UUID
	Email             string
	PasswordHash      string       `db:"password_hash"`
	LoginBlockedUntil sql.NullTime `db:"login_blocked_until"`
}

func (r *pg) GetUserByEmail(email string) (*GetUserByEmailOutput, error) {
	var output GetUserByEmailOutput
	err := r.db.QueryOne(&output, "SELECT id, email, password_hash, login_blocked_until FROM users WHERE LOWER(email) = LOWER($1) LIMIT 1", email)
	if err != nil {
		return nil, err
	}
	return &output, nil
}
