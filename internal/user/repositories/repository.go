package user

import "github.com/vagnercardosoweb/go-rest-api/sqlc/store"

type RepositoryInterface interface {
	GetUserByEmailToLogin(email string) (*store.GetUserByEmailToLoginRow, error)
}
