package user

import (
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
)

type RepositoryInterface interface {
	GetByEmail(email string) (*GetByEmailOutput, error)
}

type Repository struct {
	pgClient *postgres.Client
}

func NewRepository(db *postgres.Client) RepositoryInterface {
	return &Repository{pgClient: db}
}
