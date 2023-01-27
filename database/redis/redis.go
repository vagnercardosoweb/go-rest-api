package redis

import (
	"context"
	"rest-api/shared"

	libRedis "github.com/go-redis/redis/v9"
)

type Connection struct {
	ctx    context.Context
	client *libRedis.Client
	logger *shared.Logger
}

func NewConnection(ctx context.Context) *Connection {
	client := libRedis.NewClient(NewConfig())
	return &Connection{
		ctx:    ctx,
		client: client,
		logger: shared.NewLogger(shared.Logger{Id: "REDIS"}),
	}
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

func (c *Connection) GetClient() *libRedis.Client {
	return c.client
}
