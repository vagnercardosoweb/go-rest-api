package user

import (
	"context"

	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/sqlc/store"
)

type RepositoryPg struct {
	db *postgres.Connection
}

func NewRepositoryPg(db *postgres.Connection) RepositoryInterface {
	return &RepositoryPg{db: db}
}

func (repo *RepositoryPg) GetUserByEmailToLogin(email string) (*store.GetUserByEmailToLoginRow, error) {
	return repo.db.Queries.GetUserByEmailToLogin(context.Background(), email)
}
