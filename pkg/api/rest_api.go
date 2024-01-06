package api

import (
	"context"
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/handlers"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/middlewares"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/utils"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func New(ctx context.Context, logger *logger.Logger) *RestApi {
	restApi := &RestApi{
		logger:                 logger,
		dependencyOnTheRequest: make(map[string]any),
		appEnv:                 env.AppLocal,
		port:                   "3000",
		ctx:                    ctx,
	}

	restApi.makeHandlers()

	return restApi
}

func (r *RestApi) WithAppEnv(appEnv string) *RestApi {
	r.appEnv = appEnv
	return r
}

func (r *RestApi) WithPort(port string) *RestApi {
	r.port = port
	return r
}

func (r *RestApi) WithShutdownTimeout(timeout float64) *RestApi {
	r.shutdownTimeout = time.Duration(timeout) * time.Second
	return r
}

func (r *RestApi) WithRequestDependency(key string, value any) *RestApi {
	r.dependencyOnTheRequest[key] = value
	return r
}

func (r *RestApi) AddHandler(method string, path string, handlers ...any) *RestApi {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))

	for i, handler := range handlers {
		switch handler.(type) {
		case func(*gin.Context):
			ginHandlers[i] = handler.(func(*gin.Context))
		case func(*gin.Context) interface{}:
			ginHandlers[i] = utils.WrapperHandler(handler.(func(*gin.Context) interface{}))
		default:
			panic(errors.New(errors.Input{
				Message:   `Invalid handler "%s" for route "%s %s"`,
				Arguments: []any{handler, method, path},
			}))
		}
	}

	r.gin.Handle(method, path, ginHandlers...)
	return r
}

func (r *RestApi) Run() {
	go r.listen()
	r.logger.Info("server started on port %s", r.port)
	r.shutdown()
}

func (r *RestApi) listen() {
	r.server = &http.Server{Addr: fmt.Sprintf(":%s", r.port), Handler: r.gin}

	if err := r.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		r.logger.AddMetadata("error", err).Error("server listen error")
		os.Exit(1)
	}
}

func (r *RestApi) shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	r.logger.Info("shutdown server")

	ctx, cancel := context.WithTimeout(r.ctx, r.shutdownTimeout)
	defer cancel()

	if err := r.server.Shutdown(ctx); err != nil {
		r.logger.AddMetadata("error", err).Error("server forced shutdown")
		os.Exit(1)
	}

	select {
	case <-ctx.Done():
		r.logger.Info(`timeout of "%.0f" seconds.`, r.shutdownTimeout.Seconds())
	}

	r.logger.Info("server exiting")
	os.Exit(0)
}

func (r *RestApi) makeHandlers() {
	if r.appEnv == env.AppTest {
		gin.SetMode(gin.TestMode)
	} else if r.appEnv != env.AppLocal {
		gin.SetMode(gin.ReleaseMode)
	}

	r.gin = gin.New()

	r.gin.RedirectTrailingSlash = true
	r.gin.RemoveExtraSlash = true

	r.gin.Use(middlewares.ResponseTime)
	r.gin.Use(gzip.Gzip(gzip.BestSpeed))

	r.gin.Use(middlewares.Cors)
	r.gin.Use(middlewares.ProtectedHeaders)
	r.gin.Use(middlewares.ValidateTranslator)

	r.gin.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(r.ctx)

		for key, value := range r.dependencyOnTheRequest {
			c.Set(key, value)
		}

		c.Next()
	})

	r.gin.Use(middlewares.RequestId)
	r.gin.Use(middlewares.RequestLog)
	r.gin.Use(middlewares.ExtractAuthToken)
	r.gin.Use(gin.CustomRecovery(middlewares.Recovery))
	r.gin.Use(middlewares.ResponseError)

	r.gin.GET("/healthy", utils.WrapperHandler(handlers.Healthy))
	r.gin.GET("/favicon.ico", handlers.Favicon)

	r.gin.NoMethod(handlers.NotAllowed)
	r.gin.NoRoute(handlers.NotFound)
}
