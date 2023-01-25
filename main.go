package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"rest-api/config"
	"rest-api/handlers"
	"rest-api/shared"
)

var (
	ctx    context.Context
	logger *shared.Logger
)

func init() {
	ctx = context.Background()
	logger = shared.GetLogger()
	shared.EnvLoadFromFile()

}

func handleSignals(httpServer *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	logger.Error("Shutting down server")

	// ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	// defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		logger.AddMetadata("error", err).Error("Server forced to shutdown")
		os.Exit(1)
	}

	// select {
	// case <-ctx.Done():
	// 	logger.Warning("Timeout of 5 seconds.")
	// }

	logger.Error("Server exiting")
}

func startServer() *http.Server {
	router := gin.New()

	log.Println(shared.JwtGenerateBySubject("aaa"))

	router.Use(gin.Recovery())
	router.Use(handlers.Cors)
	router.Use(handlers.RequestId)
	router.Use(handlers.Logger)
	router.Use(handlers.Error)
	router.NoRoute(handlers.NotFound)
	router.NoMethod(handlers.NotAllowed)
	router.GET("/", handlers.Healthy)

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", shared.EnvGetByName("PORT", "3333")),
		Handler: router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.AddMetadata("error", err).Error("http server error")
			os.Exit(1)
		}
	}()

	logger.Info(
		fmt.Sprintf(
			"Server running on http://0.0.0.0:%s",
			shared.EnvGetByName("LOCAL_PORT", "3301"),
		),
	)

	// Press Cmd+C / Ctrl+C to stop.

	return httpServer
}

func main() {
	if config.IsDebug {
		shared.StartProfiler()
	}
	httpServer := startServer()
	handleSignals(httpServer)
}
