package redis

import (
	"context"
	"fmt"
	"rest-api/shared"
	"strconv"

	libRedis "github.com/go-redis/redis/v9"
)

var (
	addr = fmt.Sprintf(
		"%s:%s",
		shared.EnvRequiredByName("REDIS_HOST"),
		shared.EnvRequiredByName("REDIS_PORT"),
	)
	database = shared.EnvRequiredByName("REDIS_DATABASE")
	password = shared.EnvRequiredByName("REDIS_PASSWORD")
)

func NewRedisClient(ctx context.Context) *libRedis.Client {
	redisDb, _ := strconv.Atoi(database)
	options := &libRedis.Options{
		Addr:     addr,
		Password: password,
		DB:       redisDb,
	}
	client := libRedis.NewClient(options)
	return client
}
