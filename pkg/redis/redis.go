package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

func NewClient(
	ctx context.Context,
	options *redis.Options,
) *Client {
	client := redis.NewClient(options)
	connection := &Client{ctx: ctx, redis: client}

	err := connection.Ping()
	if err != nil {
		panic(err)
	}

	return connection
}

func FromEnv(ctx context.Context) *Client {
	return NewClient(ctx, &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", env.Required("REDIS_HOST"), env.Required("REDIS_PORT")),
		Password: env.GetAsString("REDIS_PASSWORD", ""),
		Username: env.GetAsString("REDIS_USERNAME", ""),
		DB:       env.GetAsInt("REDIS_DATABASE", "0"),
	})
}

func (c *Client) Get(key string, dest any) error {
	if reflect.ValueOf(dest).Kind() != reflect.Ptr {
		return fmt.Errorf("Redis#Get('%s') dest must be pointer", key)
	}

	bytes, err := c.redis.Get(c.ctx, key).Bytes()

	if errors.Is(err, redis.Nil) {
		return nil
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, dest)
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

func (c *Client) FlushAll() error {
	return c.redis.FlushAll(c.ctx).Err()
}

func (c *Client) Close() error {
	return c.redis.Close()
}
