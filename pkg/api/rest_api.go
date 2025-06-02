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

	api := &Api{
		ctx:         ctx,
		logger:      logger,
		environment: env.Development,
		values:      make(map[string]any),
		onStart:     make([]func(api *Api), 0),
		onShutdown:  make([]func(api *Api, code string), 0),
		routes:      make([]*Route, 0),
		server: &http.Server{
			ReadTimeout:       30 * time.Second,
			MaxHeaderBytes:    2 << 20, // 2 MB
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      30 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}

	api.WithPort(env.GetAsString("PORT", "3000"))
	api.WithShutdownTimeout(env.GetAsFloat64("SHUTDOWN_TIMEOUT", "0"))

	return api
}

func (api *Api) Logger() *logger.Logger {
	return api.logger
}

func (api *Api) GetServer() *http.Server {
	return api.server
}

func (api *Api) WithEnv(env string) *Api {
	api.environment = env
	return api
}

func (api *Api) WithPort(port string) *Api {
	api.server.Addr = fmt.Sprintf(":%s", port)
	return api
}

func (api *Api) WithShutdownTimeout(timeout float64) *Api {
	api.shutdownTimeout = time.Duration(timeout) * time.Second
	return api
}

func (api *Api) WithValue(key string, value any) *Api {
	api.values[key] = value

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

func (api *Api) OnStart(fn func(api *Api)) *Api {
	api.onStart = append(api.onStart, fn)
	return api
}

func (api *Api) OnShutdown(fn func(api *Api, code string)) *Api {
	api.onShutdown = append(api.onShutdown, fn)
	return api
}

func (api *Api) Run() {
	api.Start()

	go func() {
		if err := api.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			api.logger.AddField("error", err).Error("server listen error")
			os.Exit(1)
		}
	}()

	api.logger.Info(`server is running on port "%s"`, api.server.Addr)
	api.shutdown()
}

func (api *Api) Start() {
	api.setupGin()
	api.setupHandlers()

	// Run onStart callbacks
	for _, fn := range api.onStart {
		fn(api)
	}
}

func (api *Api) setupHandlers() {
	for _, route := range api.routes {
		handlers := make([]gin.HandlerFunc, len(api.routes))

		for i, handler := range route.Handlers {
			switch h := handler.(type) {
			case func(*gin.Context):
				handlers[i] = h
			case func(*gin.Context) interface{}:
				handlers[i] = apiresponse.Wrapper(h)
			default:
				panic(fmt.Errorf(
					`invalid handler "%s" for route "%s %s"`,
					handler, route.Method, route.Path,
				))
			}
		}

		api.gin.Handle(
			route.Method,
			route.Path,
			handlers...,
		)
	}
}

func (api *Api) shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	code := <-quit
	api.logger.Info(`server exited with code "%d"`, code)

	// Run onShutdown callbacks
	for _, fn := range api.onShutdown {
		fn(api, code.String())
	}

	ctx, cancel := context.WithTimeout(api.ctx, api.shutdownTimeout)
	defer cancel()

	if err := api.server.Shutdown(ctx); err != nil {
		api.logger.
			AddField("error", err).
			Error("shutdown server error")
	}

	// Wait for the context to be done
	<-ctx.Done()

	api.logger.Info(`server exiting with timeout of "%s"`, api.shutdownTimeout.String())

	os.Exit(0)
}

func (api *Api) setupGin() {
	if api.environment == env.Test {
		gin.SetMode(gin.TestMode)
	} else if api.environment == env.Production {
		gin.SetMode(gin.ReleaseMode)
	}

	api.gin = gin.New()
	api.server.Handler = api.gin

	api.gin.RedirectTrailingSlash = true
	api.gin.RemoveExtraSlash = true

	api.gin.Use(middlewares.ResponseTime)
	api.gin.Use(gzip.Gzip(gzip.BestSpeed))

	api.gin.Use(middlewares.Cors)
	api.gin.Use(middlewares.RequestId)
	api.gin.Use(middlewares.SecurityHeaders)
	api.gin.Use(middlewares.NoCacheHeaders)
	api.gin.Use(middlewares.BearerToken)

	api.gin.Use(func(c *gin.Context) {
		c.Request = c.Request.WithContext(api.ctx)

		// Response always JSON
		c.Header("Content-Type", "application/json; charset=utf-8")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusOK)
			return
		}

		requestId := apicontext.RequestId(c)
		c.Set(logger.CtxKey, api.logger.WithId(requestId))

		for key, value := range api.values {
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

	api.gin.GET("/healthy", handlers.Healthy)
	api.gin.GET("/favicon.ico", handlers.Favicon)
	api.gin.GET("/timestamp", handlers.Timestamp)

	api.gin.NoMethod(handlers.NotAllowed)
	api.gin.NoRoute(handlers.NotFound)
}
