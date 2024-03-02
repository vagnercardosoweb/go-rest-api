package tests

import (
	"github.com/vagnercardosoweb/go-rest-api/internal/events"
	"github.com/vagnercardosoweb/go-rest-api/internal/handlers/user"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

type RestApiSuite struct {
	ContainerTestSuite
	RestApi *api.RestApi
}

func (r *RestApiSuite) SetupSuite() {
	r.ContainerTestSuite.SetupSuite()

	r.RestApi = api.New(r.Ctx, r.Logger)
	r.RestApi.WithAppEnv(env.AppTest)

	r.RestApi.WithValue(redis.CtxKey, r.RedisClient)
	r.RestApi.WithValue(postgres.CtxKey, r.PgClient)
	r.RestApi.WithValue(token.ClientCtxKey, token.NewJwtFromEnv())
	r.RestApi.WithValue(events.CtxKey, events.New(r.PgClient, r.RedisClient))
	r.RestApi.WithValue(logger.CtxKey, r.Logger)

	user.MakeHandlers(r.RestApi)

	r.RestApi.Start()
}

func (r *RestApiSuite) TearDownSuite() {
	r.ContainerTestSuite.TearDownSuite()
}
