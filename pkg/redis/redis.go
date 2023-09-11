package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/redis/v9"
)

type Client struct {
	redis *redis.Client
	ctx   context.Context
}

func NewClient(ctx context.Context) *Client {
	client := redis.NewClient(newConfig())
	connection := &Client{
		ctx:   ctx,
		redis: client,
	}
	err := connection.Ping()
	if err != nil {
		panic(err)
	}
	return connection
}

func (c *Client) Get(key string, dest any) error {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return fmt.Errorf("Redis#Get('%s') dest must be pointer", key)
	}
	result, err := c.redis.Get(c.ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), dest)
}

func (c *Client) Set(key string, value any, expiration time.Duration) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.redis.Set(
		c.ctx,
		key,
		jsonValue,
		expiration,
	).Err()
}

func (c *Client) Has(key string) (bool, error) {
	cmd := c.redis.Exists(c.ctx, key)
	return c.checkResultCmd(cmd)
}

func (c *Client) checkResultCmd(cmd *redis.IntCmd) (bool, error) {
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	return cmd.Val() > 0, nil
}

func (c *Client) Del(key string) (bool, error) {
	cmd := c.redis.Del(c.ctx, key)
	return c.checkResultCmd(cmd)
}

func (c *Client) Ping() error {
	result := c.redis.Ping(c.ctx)
	return result.Err()
}

func (c *Client) Close() error {
	return c.redis.Close()
}
