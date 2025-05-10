package main

import (
	"context"

	"github.com/vagnercardosoweb/go-rest-api/internal/events"
	"github.com/vagnercardosoweb/go-rest-api/internal/handlers/user"
	"github.com/vagnercardosoweb/go-rest-api/internal/schedules"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api"
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

	restApi := api.
		New(ctx, appLogger).
		WithAppEnv(env.GetAppEnv()).
		WithShutdownTimeout(env.GetAsFloat64("SHUTDOWN_TIMEOUT", "0")).
		WithPort(env.Required("PORT")).
		WithValue(token.ClientCtxKey, token.NewJwtFromEnv()).
		WithValue(events.CtxKey, events.New(pgClient, redisClient)).
		WithValue(redis.CtxKey, redisClient).
		WithValue(postgres.CtxKey, pgClient).
		WithValue(logger.CtxKey, appLogger)

	// Setup handlers
	user.MakeHandlers(restApi)

	if env.IsSchedulerEnabled() {
		go schedules.New(
			pgClient,
			redisClient,
			appLogger.WithId("SCHEDULER"),
		).Run()
	}

	if env.GetAsBool("DEBUG") {
		monitoring.RunProfiler(appLogger)
	}

	restApi.Run()
}
