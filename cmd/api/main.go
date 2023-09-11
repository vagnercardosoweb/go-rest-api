package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"

	"github.com/gin-contrib/gzip"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"

	"github.com/gin-gonic/gin"

	"github.com/vagnercardosoweb/go-rest-api/cmd/api/middlewares"
	"github.com/vagnercardosoweb/go-rest-api/cmd/api/routes"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/monitoring"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
)

var (
	ctx           context.Context
	httpServer    *http.Server
	pgClient      *postgres.Client
	redisClient   *redis.Client
	tokenInstance token.Token
	appLogger     *logger.Logger
)

func init() {
	env.Load()
	ctx = context.Background()
	appLogger = logger.New()

	tokenInstance = token.NewJwt([]byte(env.Required("JWT_SECRET_KEY")), config.GetExpiresInJwt())
	ctx = context.WithValue(ctx, config.TokenCtxKey, tokenInstance)

	pgClient = postgres.NewClient(ctx, appLogger, postgres.Default)
	ctx = context.WithValue(ctx, config.PgClientCtxKey, pgClient)

	redisClient = redis.NewClient(ctx)
	ctx = context.WithValue(ctx, config.RedisClientCtxKey, redisClient)

	httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", env.Required("PORT")),
		Handler: handler(),
	}
}

func shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	appLogger.Info("Shutting down server")

	timeout := config.GetShutdownTimeout() * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown: %v", err.Error())
		os.Exit(1)
	}

	<-ctx.Done()

	appLogger.Info("Server exiting of %v seconds.", timeout)
}

func handler() *gin.Engine {
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.RemoveExtraSlash = true
	router.RedirectTrailingSlash = true

	router.Use(gzip.Gzip(gzip.BestSpeed))
	router.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("Content-Type", "application/json")

		c.Set(config.PgClientCtxKey, pgClient)
		c.Set(config.RedisClientCtxKey, redisClient)
		c.Set(config.RequestLoggerCtxKey, appLogger)
		c.Set(config.TokenCtxKey, tokenInstance)

		c.Next()
	})

	middlewares.Setup(router)
	routes.Setup(router)

	return router
}

func main() {
	defer pgClient.Close()
	defer redisClient.Close()

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appLogger.Error("Server listen error: %v", err.Error())
			os.Exit(1)
		}
	}()

	appLogger.Info(
		"Server running on http://0.0.0.0%s",
		httpServer.Addr,
	)

	if config.IsDebug() {
		monitoring.RunProfiler(appLogger)
	}

	shutdown()
}
