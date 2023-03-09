package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v9"
)

type Connection struct {
	ctx    context.Context
	client *redis.Client
}

func Connect(ctx context.Context) *Connection {
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

func (c *Connection) Get(key string, dest any) error {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return fmt.Errorf("Redis#Get('%s') dest must be pointer", key)
	}
	result, err := c.client.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), dest)
}

func (c *Connection) Set(key string, value any, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.client.Set(
		c.ctx,
		key,
		jsonValue,
		expiration,
	).Err()
}

func (c *Connection) Has(key string) (bool, error) {
	cmd := c.client.Exists(c.ctx, key)
	return c.checkResultCmd(cmd)
}

func (c *Connection) checkResultCmd(cmd *redis.IntCmd) (bool, error) {
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	return cmd.Val() > 0, nil
}

func (c *Connection) Del(key string) (bool, error) {
	cmd := c.client.Del(c.ctx, key)
	return c.checkResultCmd(cmd)
}

func (c *Connection) Ping() error {
	result := c.client.Ping(c.ctx)
	return result.Err()
}

func (c *Connection) Close() error {
	return c.client.Close()
}
