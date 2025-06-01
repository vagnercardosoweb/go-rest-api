package main

import (
	"context"
	"fmt"

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
	"github.com/vagnercardosoweb/go-rest-api/pkg/slack"
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

	eventManager := events.NewManager(pgClient, redisClient)
	restApi := api.New(ctx, appLogger).
		WithEnv(env.GetAppEnv()).
		WithValue(redis.CtxKey, redisClient).
		WithValue(token.CtxClientKey, token.JwtFromEnv()).
		WithValue(password.CtxKey, password.NewBcrypt()).
		WithValue(postgres.CtxKey, func(c *gin.Context) any {
			return pgClient.WithLogger(apicontext.Logger(c))
		}).
		WithValue(events.CtxKey, func(c *gin.Context) any {
			return eventManager.WithLogger(apicontext.Logger(c))
		}).
		OnStart(func(api *api.Api) {
			go func() {
				if env.IsAlertOnServerStart() {
					_ = slack.NewAlert().
						AddField("message", fmt.Sprintf(`server is running on port "%s"`, api.GetServer().Addr), false).
						Send()
				}
			}()
		}).
		OnShutdown(func(api *api.Api, code string) {
			go func() {
				if env.IsAlertOnServerClose() {
					_ = slack.NewAlert().
						WithColor(slack.ColorError).
						AddField("message", fmt.Sprintf(`server exited with code "%s"`, code), false).
						Send()
				}
			}()
		})

	// Make handlers
	user.MakeHandlers(restApi)

	if env.IsSchedulerEnabled() {
		go schedules.New(pgClient, redisClient, appLogger).Run()
	}

	if env.GetAsBool("PROFILER_ENABLED") {
		monitoring.RunProfiler(appLogger)
	}

	restApi.Run()
}
