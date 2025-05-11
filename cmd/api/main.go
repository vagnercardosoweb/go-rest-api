package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/internal/events"
	"github.com/vagnercardosoweb/go-rest-api/internal/handlers/user"
	"github.com/vagnercardosoweb/go-rest-api/internal/schedules"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/monitoring"
	"github.com/vagnercardosoweb/go-rest-api/pkg/password"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

func main() {
	env.Load()

	ctx := context.Background()
	appLogger := logger.New()

	redisClient := redis.FromEnv(ctx)
	defer redisClient.Close()

	pgClient := postgres.FromEnv(ctx, appLogger)
	defer pgClient.Close()

	restApi := api.New(ctx, appLogger)

	restApi.WithEnv(env.GetAppEnv())
	restApi.WithPort(env.Required("PORT"))

	restApi.WithShutdownTimeout(env.GetAsFloat64("SHUTDOWN_TIMEOUT", "0"))

	restApi.WithValue(token.CtxClientKey, token.JwtFromEnv())
	restApi.WithValue(password.CtxKey, password.NewBcrypt())
	restApi.WithValue(redis.CtxKey, redisClient)

	restApi.WithValue(postgres.CtxKey, func(c *gin.Context) any {
		return pgClient.WithLogger(apicontext.Logger(c))
	})

	eventManager := events.NewManager(pgClient, redisClient)
	restApi.WithValue(events.CtxKey, func(c *gin.Context) any {
		return eventManager.WithLogger(apicontext.Logger(c))
	})

	user.MakeHandlers(restApi)

	if env.IsSchedulerEnabled() {
		go schedules.New(pgClient, redisClient, appLogger).Run()
	}

	if env.GetAsBool("PROFILER_ENABLED") {
		monitoring.RunProfiler(appLogger)
	}

	restApi.Run()
}
