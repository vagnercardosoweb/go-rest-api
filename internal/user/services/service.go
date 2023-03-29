package user

import (
	user "github.com/vagnercardosoweb/go-rest-api/internal/user/repositories"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password_hash"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

type ServiceInterface interface {
	Login(email, password string) (string, error)
}

type Service struct {
	passwordHash   password_hash.PasswordHash
	userRepository user.RepositoryInterface
	token          token.Token
}

func NewService(userRepository user.RepositoryInterface, passwordHash password_hash.PasswordHash, token token.Token) ServiceInterface {
	return &Service{userRepository: userRepository, passwordHash: passwordHash, token: token}
}
