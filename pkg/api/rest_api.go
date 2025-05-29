package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	apicontext "github.com/vagnercardosoweb/go-rest-api/pkg/api/context"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/handlers"
	"github.com/vagnercardosoweb/go-rest-api/pkg/api/middlewares"
	apiresponse "github.com/vagnercardosoweb/go-rest-api/pkg/api/response"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

func New(ctx context.Context, logger *logger.Logger) *Api {
	return &Api{
		ctx:          ctx,
		port:         "3001",
		environment:  env.Development,
		dependencies: make(map[string]any),
		routes:       make([]*Route, 0),
		logger:       logger,
	}
}

func (api *Api) Logger() *logger.Logger {
	return api.logger
}

func (api *Api) WithEnv(env string) *Api {
	api.environment = env
	return api
}

func (api *Api) WithPort(port string) *Api {
	api.port = port
	return api
}

func (api *Api) WithShutdownTimeout(timeout float64) *Api {
	api.shutdownTimeout = time.Duration(timeout) * time.Second
	return api
}

func (api *Api) WithValue(key string, value any) *Api {
	api.dependencies[key] = value

	type keyType string
	api.ctx = context.WithValue(api.ctx, keyType(key), value)

	return api
}

func (api *Api) Get(path string, handlers ...any) *Api {
	return api.AddHandler(http.MethodGet, path, handlers...)
}

func (api *Api) Post(path string, handlers ...any) *Api {
	return api.AddHandler(http.MethodPost, path, handlers...)
}

func (api *Api) Put(path string, handlers ...any) *Api {
	return api.AddHandler(http.MethodPut, path, handlers...)
}

func (api *Api) Patch(path string, handlers ...any) *Api {
	return api.AddHandler(http.MethodPatch, path, handlers...)
}

func (api *Api) Delete(path string, handlers ...any) *Api {
	return api.AddHandler(http.MethodDelete, path, handlers...)
}

func (api *Api) Group(path string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return api.gin.Group(path, handlers...)
}

func (api *Api) AddHandler(method string, path string, handlers ...any) *Api {
	api.routes = append(api.routes, &Route{method, handlers, path})
	return api
}

func (api *Api) TestRequest(request *http.Request) *httptest.ResponseRecorder {
	if request.Method == http.MethodPost || request.Method == http.MethodPut {
		request.Header.Set("Content-Type", "application/json")
	}

	rr := httptest.NewRecorder()
	api.gin.ServeHTTP(rr, request)

	return rr
}

func (api *Api) Run() {
	api.Start()
	go api.listen()
	api.logger.Info("server started on port %s", api.port)
	api.shutdown()
}

func (api *Api) Start() {
	api.setupGin()
	api.makeRoutes()
}

func (api *Api) makeRoutes() {
	for _, route := range api.routes {
		handlers := make([]gin.HandlerFunc, len(api.routes))

		for i, handler := range route.Handlers {
			switch h := handler.(type) {
			case func(*gin.Context):
				handlers[i] = h
			case func(*gin.Context) interface{}:
				handlers[i] = apiresponse.Wrapper(h)
			default:
				panic(errors.New(errors.Input{
					Message:   `Invalid handler "%s" for route "%s %s"`,
					Arguments: []any{handler, route.Method, route.Path},
				}))
			}
		}

		api.gin.Handle(
			route.Method,
			route.Path,
			handlers...,
		)
	}
}

func (api *Api) listen() {
	api.server = &http.Server{Addr: fmt.Sprintf(":%s", api.port), Handler: api.gin}

	if err := api.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		api.logger.AddMetadata("error", err).Error("server listen error")
		os.Exit(1)
	}
}

func (api *Api) shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	api.logger.Info("shutdown server")

	ctx, cancel := context.WithTimeout(api.ctx, api.shutdownTimeout)
	defer cancel()

	if err := api.server.Shutdown(ctx); err != nil {
		api.logger.AddMetadata("error", err).Error("server forced shutdown")
		os.Exit(1)
	}

	<-ctx.Done()
	api.logger.Info(`timeout of "%.0f" seconds.`, api.shutdownTimeout.Seconds())

	api.logger.Info("server exiting")
	os.Exit(0)
}

func (api *Api) setupGin() {
	if api.environment == env.Test {
		gin.SetMode(gin.TestMode)
	} else if api.environment == env.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	api.gin = gin.New()

	api.gin.RedirectTrailingSlash = true
	api.gin.RemoveExtraSlash = true

	api.gin.Use(middlewares.ResponseTime)
	api.gin.Use(gzip.Gzip(gzip.BestSpeed))

	api.gin.Use(middlewares.Cors)
	api.gin.Use(middlewares.RequestId)
	api.gin.Use(middlewares.BearerToken)
	api.gin.Use(middlewares.Headers)

	api.gin.Use(func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		c.Request = c.Request.WithContext(api.ctx)

		requestId := apicontext.RequestId(c)
		c.Set(logger.CtxKey, api.logger.WithId(requestId))

		for key, value := range api.dependencies {
			if handler, ok := value.(func(*gin.Context) any); ok {
				value = handler(c)
			}

			c.Set(key, value)
		}

		c.Next()
	})

	api.gin.Use(middlewares.RequestLog)
	api.gin.Use(middlewares.Translator)

	api.gin.Use(gin.CustomRecovery(middlewares.Recovery))
	api.gin.Use(middlewares.ResponseError)

	api.gin.GET("/healthy", apiresponse.Wrapper(handlers.Healthy))
	api.gin.GET("/favicon.ico", handlers.Favicon)
	api.gin.GET("/timestamp", handlers.Timestamp)

	api.gin.NoMethod(handlers.NotAllowed)
	api.gin.NoRoute(handlers.NotFound)
}
