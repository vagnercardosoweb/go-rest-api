package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	redis *redis.Client
	ctx   context.Context
}
