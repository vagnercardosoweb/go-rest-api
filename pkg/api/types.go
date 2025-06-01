package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

type Route struct {
	Method   string
	Handlers []any
	Path     string
}

type Api struct {
	ctx             context.Context
	logger          *logger.Logger
	routes          []*Route
	environment     string
	shutdownTimeout time.Duration
	onStart         []func(api *Api)
	onShutdown      []func(api *Api, code string)
	values          map[string]any
	server          *http.Server
	gin             *gin.Engine
}
