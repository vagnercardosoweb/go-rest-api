package user

import (
	"github.com/vagnercardosoweb/go-rest-api/pkg/api"
)

func MakeHandlers(restApi *api.Api) {
	restApi.Post("/login", Login)
}
