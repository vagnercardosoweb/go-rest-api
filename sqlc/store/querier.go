// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0

package store

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateUser(ctx context.Context, arg *CreateUserParams) error
	GetUserByEmailToLogin(ctx context.Context, email string) (*GetUserByEmailToLoginRow, error)
	GetUsers(ctx context.Context, limit int32) ([]*GetUsersRow, error)
	GetWalletById(ctx context.Context, id uuid.UUID) (*GetWalletByIdRow, error)
	GetWallets(ctx context.Context, limit int32) ([]*GetWalletsRow, error)
}

var _ Querier = (*Queries)(nil)
