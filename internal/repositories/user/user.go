package user

import (
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
)

type instance struct {
	pgClient *postgres.Client
}

type Repository interface {
	GetByEmail(email string) (*GetByEmailOutput, error)
	Create(input *CreateInput) (*CreateOutput, error)
}

func New(pgClient *postgres.Client) Repository {
	return &instance{pgClient: pgClient}
}
