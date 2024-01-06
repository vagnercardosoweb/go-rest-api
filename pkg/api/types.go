package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"net/http"
	"time"
)

type Route struct {
	Method   string
	Handlers []any
	Path     string
}

type RestApi struct {
	ctx                    context.Context
	logger                 *logger.Logger
	port                   string
	shutdownTimeout        time.Duration
	dependencyOnTheRequest map[string]any
	server                 *http.Server
	appEnv                 string
	gin                    *gin.Engine
}
