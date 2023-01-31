package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/api/middlewares"
	"github.com/vagnercardosoweb/go-rest-api/api/routes"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/database/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/database/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/monitoring"
)

var (
	ctx             context.Context
	httpServer      *http.Server
	dbConnection    *postgres.Connection
	redisConnection *redis.Connection
)

func init() {
	env.LoadFromFile()

	ctx = context.Background()

	dbConnection = postgres.NewConnection(ctx)
	ctx = context.WithValue(ctx, config.DbConnectionContextKey, dbConnection)

	redisConnection = redis.NewConnection(ctx)
	ctx = context.WithValue(ctx, config.RedisConnectionContextKey, redisConnection)

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

	select {
	case <-ctx.Done():
		logger.Warn("Timeout shutdown of %v seconds.", timeout)
	}

	logger.Error("Server exiting")
}

func handler() *gin.Engine {
	if config.IsProduction {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.RemoveExtraSlash = true
	router.RedirectTrailingSlash = true

	router.Use(func(c *gin.Context) {
		c.Set(config.DbConnectionContextKey, ctx.Value(config.DbConnectionContextKey))
		c.Set(config.RedisConnectionContextKey, ctx.Value(config.RedisConnectionContextKey))
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	})

	router.Use(gin.Recovery())

	middlewares.Setup(router)
	routes.Setup(router)

	return router
}

func main() {
	defer dbConnection.Close()
	defer redisConnection.Close()

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server listen error: %v", err.Error())
			os.Exit(1)
		}
	}()

	logger.Info(
		fmt.Sprintf(
			"Server running on http://0.0.0.0:%s",
			env.Get("LOCAL_PORT", "3301"),
		),
	)

	if config.IsDebug {
		monitoring.RunProfiler()
	}

	shutdown()
}
