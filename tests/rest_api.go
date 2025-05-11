package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/internal/events"
	"github.com/vagnercardosoweb/go-rest-api/internal/handlers/user"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

type RestApiSuite struct {
	ContainerTestSuite
	RestApi *api.Api
}

func (r *RestApiSuite) SetupSuite() {
	r.ContainerTestSuite.SetupSuite()

	r.RestApi = api.New(r.Ctx, r.Logger)
	r.RestApi.WithEnv(env.Test)

	r.RestApi.WithValue(password.CtxKey, password.NewBcrypt())
	r.RestApi.WithValue(token.CtxClientKey, token.JwtFromEnv())
	r.RestApi.WithValue(redis.CtxKey, r.RedisClient)

	r.RestApi.WithValue(postgres.CtxKey, func(c *gin.Context) any {
		return r.PgClient.WithLogger(apicontext.Logger(c))
	})

	eventManager := events.NewManager(r.PgClient, r.RedisClient)
	r.RestApi.WithValue(events.CtxKey, func(c *gin.Context) any {
		return eventManager.WithLogger(apicontext.Logger(c))
	})

	user.MakeHandlers(r.RestApi)

	r.RestApi.Start()
}

func (r *RestApiSuite) TearDownSuite() {
	r.ContainerTestSuite.TearDownSuite()
}
