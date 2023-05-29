package main

import (
	"context"
	"fmt"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/vagnercardosoweb/go-rest-api/cmd/api/middlewares"
	"github.com/vagnercardosoweb/go-rest-api/cmd/api/routes"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/monitoring"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
)

var (
	ctx          context.Context
	httpServer   *http.Server
	postgresConn *postgres.Connection
	redisConn    *redis.Connection
	appLogger    *logger.Logger
)

func init() {
	env.LoadFromLocal()
	ctx = context.Background()
	appLogger = logger.New()

	postgresConn = postgres.Connect(ctx, postgres.Default)
	ctx = context.WithValue(ctx, config.PgConnectCtxKey, postgresConn)

	redisConn = redis.Connect(ctx)
	ctx = context.WithValue(ctx, config.RedisConnectCtxKey, redisConn)

	httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", env.Get("PORT", "3333")),
		Handler: handler(),
	}
}

func shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	appLogger.Error("Shutting down server")

	timeout := config.GetShutdownTimeout() * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		appLogger.Error("Server forced to shutdown: %v", err.Error())
		os.Exit(1)
	}

	<-ctx.Done()

	appLogger.Error("Server exiting of %v seconds.", timeout)
}

func handler() *gin.Engine {
	if config.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.RemoveExtraSlash = true
	router.RedirectTrailingSlash = true

	router.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(ctx)

		c.Set(config.RequestLoggerCtxKey, appLogger)
		c.Set(config.PgConnectCtxKey, postgresConn)
		c.Set(config.RedisConnectCtxKey, redisConn)

		c.Next()
	})

	middlewares.Setup(router)
	routes.Setup(router)

	return router
}

func main() {
	defer redisConn.Close()
	defer postgresConn.Close()

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Error("Server listen error: %v", err.Error())
			os.Exit(1)
		}
	}()

	appLogger.Info(
		"Server running on http://0.0.0.0:%s",
		env.Get("LOCAL_PORT", "3301"),
	)

	if config.IsDebug {
		monitoring.RunProfiler(appLogger)
	}

	shutdown()
}
