package redis

import (
	"context"

	"github.com/go-redis/redis/v9"
)

type Connection struct {
	ctx    context.Context
	client *redis.Client
}

func NewConnection(ctx context.Context) *Connection {
	client := redis.NewClient(newConfig())
	connection := &Connection{
		ctx:    ctx,
		client: client,
	}
	err := connection.Ping()
	if err != nil {
		panic(err)
	}
	return connection
}

func (c *Connection) Get() {

}

func (c *Connection) Ping() error {
	result := c.client.Ping(c.ctx)
	return result.Err()
}

func (c *Connection) Close() error {
	return c.client.Close()
}

func (c *Connection) GetClient() *redis.Client {
	return c.client
}
