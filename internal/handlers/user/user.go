package user

import (
	"github.com/vagnercardosoweb/go-rest-api/pkg/api"
)

func MakeHandlers(api *api.Api) {
	api.Post("/login", Login)
}
