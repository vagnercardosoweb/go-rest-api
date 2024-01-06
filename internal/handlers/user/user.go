package user

import (
	"github.com/vagnercardosoweb/go-rest-api/pkg/api"
	"net/http"
)

func MakeHandlers(restApi *api.RestApi) {
	restApi.AddHandler(http.MethodPost, "/login", Login)
}
