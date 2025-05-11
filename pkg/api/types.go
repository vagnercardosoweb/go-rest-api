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
	port            string
	routes          []*Route
	shutdownTimeout time.Duration
	dependencies    map[string]any
	server          *http.Server
	environment     string
	gin             *gin.Engine
}
