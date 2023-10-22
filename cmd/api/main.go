package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/vagnercardosoweb/go-rest-api/cmd/api/middlewares"
	"github.com/vagnercardosoweb/go-rest-api/cmd/api/routes"
	"github.com/vagnercardosoweb/go-rest-api/internal/schedules"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/monitoring"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/token"
)

var (
	ctx          context.Context
	pgClient     *postgres.Client
	httpServer   *http.Server
	redisClient  *redis.Client
	tokenManager token.Token
	appLogger    *logger.Logger
)

func init() {
	env.Load()
	ctx = context.Background()

	appLogger = logger.New()
	ctx = context.WithValue(ctx, config.RequestLoggerCtxKey, appLogger)

	tokenManager = token.NewJwt(config.JwtSecretKey(), config.JwtExpiresIn())
	ctx = context.WithValue(ctx, config.TokenManagerCtxKey, tokenManager)

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

	appLogger.Info("shutdown server")

	timeout := time.Duration(env.GetInt("SHUTDOWN_TIMEOUT", "0")) * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		appLogger.AddMetadata("error", err).Error("server forced shutdown")
		os.Exit(1)
	}

	select {
	case <-ctx.Done():
		appLogger.Info(`timeout of "%.0f" seconds.`, timeout.Seconds())
	}

	appLogger.Info("server exiting")
}

func handler() *gin.Engine {
	if config.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	h := gin.New()
	h.RemoveExtraSlash = true
	h.RedirectTrailingSlash = true

	h.Use(gzip.Gzip(gzip.BestSpeed))

	h.Use(middlewares.ResponseTime)
	h.Use(middlewares.Cors)

	h.Use(func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Request = c.Request.WithContext(ctx)
		c.Writer.Header().Set("Content-Type", "application/json")

		c.Set(config.PgClientCtxKey, pgClient)
		c.Set(config.TokenManagerCtxKey, tokenManager)
		c.Set(config.RequestLoggerCtxKey, appLogger)
		c.Set(config.RedisClientCtxKey, redisClient)

		c.Next()
	})

	h.Use(middlewares.RequestId)
	h.Use(middlewares.RequestLog)
	h.Use(middlewares.ExtractAuthToken)
	h.Use(gin.CustomRecovery(middlewares.Recovery))
	h.Use(middlewares.ResponseError)

	routes.Setup(h)

	return h
}

func main() {
	defer pgClient.Close()
	defer redisClient.Close()

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			appLogger.AddMetadata("error", err).Error("server listen error")
			os.Exit(1)
		}
	}()

	appLogger.Info("server started on port %s", httpServer.Addr)

	if config.IsSchedulerEnabled() {
		scheduler := schedules.New(pgClient, redisClient, appLogger.WithID("SCHEDULER"))
		go scheduler.Run()
	}

	if config.IsDebug() {
		monitoring.RunProfiler(appLogger)
	}

	shutdown()
}
