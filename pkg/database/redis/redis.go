package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

type Connection struct {
	ctx    context.Context
	client *redis.Client
	logger *logger.Input
}

func NewConnection(ctx context.Context) *Connection {
	client := redis.NewClient(newConfig())
	return &Connection{
		ctx:    ctx,
		client: client,
		logger: logger.New(logger.Input{Id: "REDIS"}),
	}
}

func (c *Connection) Get() {

}

func (c *Connection) Ping() error {
	c.logger.Debug("Ping connection")
	result := c.client.Ping(c.ctx)
	return result.Err()
}

func (c *Connection) Close() error {
	c.logger.Debug("Closing connection")
	return c.client.Close()
}

func (c *Connection) GetClient() *redis.Client {
	return c.client
}
