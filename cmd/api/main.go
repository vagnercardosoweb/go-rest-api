package main

import (
	"context"
	"github.com/vagnercardosoweb/go-rest-api/internal/handlers/user"
	"github.com/vagnercardosoweb/go-rest-api/internal/schedules"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api"
	apiutils "github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/monitoring"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

func main() {
	env.Load()

	ctx := context.Background()
	appLogger := logger.New()

	pgClient := postgres.NewFromEnv(ctx, appLogger)
	defer pgClient.Close()

	redisClient := redis.NewFromEnv(ctx, appLogger)
	defer redisClient.Close()

	restApi := api.New(ctx, appLogger).
		WithAppEnv(env.GetAppEnv()).
		WithPort(env.Required("PORT")).
		WithShutdownTimeout(env.GetAsFloat64("SHUTDOWN_TIMEOUT", "0")).
		WithRequestDependency(apiutils.PgClientCtxKey, pgClient).
		WithRequestDependency(apiutils.RequestLoggerCtxKey, appLogger).
		WithRequestDependency(apiutils.TokenClientCtxKey, token.NewJwtFromEnv()).
		WithRequestDependency(apiutils.RedisClientCtxKey, redisClient)

	user.MakeHandlers(restApi)

	if env.IsSchedulerEnabled() {
		go schedules.New(pgClient, redisClient, appLogger).Run()
	}

	if env.GetAsBool("DEBUG") {
		monitoring.RunProfiler(appLogger)
	}

	restApi.Run()
}
