package user

import (
	"net/http"

	"github.com/vagnercardosoweb/go-rest-api/pkg/api"
)

func MakeHandlers(restApi *api.RestApi) {
	restApi.AddHandler(http.MethodPost, "/login", Login)
}
