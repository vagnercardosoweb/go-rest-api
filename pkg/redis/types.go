package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

type Client struct {
	redis  *redis.Client
	logger *logger.Logger
	ctx    context.Context
}
