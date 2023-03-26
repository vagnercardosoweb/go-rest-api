package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/sqlc/store"

	"github.com/gin-gonic/gin"

	"github.com/vagnercardosoweb/go-rest-api/cmd/api/middlewares"
	"github.com/vagnercardosoweb/go-rest-api/cmd/api/routes"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/database/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/database/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/monitoring"
)

var (
	ctx          context.Context
	httpServer   *http.Server
	postgresConn *postgres.Connection
	storeQueries *store.Queries
	redisConn    *redis.Connection
)

func init() {
	env.LoadFromLocal()
	ctx = context.Background()

	postgresConn = postgres.Connect(ctx)
	ctx = context.WithValue(ctx, config.PgConnectCtxKey, postgresConn)

	redisConn = redis.Connect(ctx)
	ctx = context.WithValue(ctx, config.RedisConnectCtxKey, redisConn)

	storeQueries = store.New(postgresConn.GetSqlx())
	ctx = context.WithValue(ctx, config.StoreQueriesCtx, storeQueries)

	httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", env.Get("PORT", "3333")),
		Handler: handler(),
	}
}

func shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	logger.Error("Shutting down server")

	timeout := config.GetShutdownTimeout() * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown: %v", err.Error())
		os.Exit(1)
	}

	<-ctx.Done()

	logger.Warn("Timeout shutdown of %v seconds.", timeout)
	logger.Error("Server exiting")
}

func handler() *gin.Engine {
	if config.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.RemoveExtraSlash = true
	router.RedirectTrailingSlash = true

	router.Use(gin.Recovery())
	router.Use(func(c *gin.Context) {
		c.Set(config.PgConnectCtxKey, ctx.Value(config.PgConnectCtxKey))
		c.Set(config.RedisConnectCtxKey, ctx.Value(config.RedisConnectCtxKey))
		c.Set(config.StoreQueriesCtx, ctx.Value(config.StoreQueriesCtx))
		c.Request = c.Request.WithContext(ctx)
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
			logger.Error("Server listen error: %v", err.Error())
			os.Exit(1)
		}
	}()

	logger.Info(
		"Server running on http://0.0.0.0:%s",
		env.Get("LOCAL_PORT", "3301"),
	)

	if config.IsDebug {
		monitoring.RunProfiler()
	}

	shutdown()
}
