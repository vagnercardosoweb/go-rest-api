package database

import (
	"context"
	"fmt"
	"rest-api/shared"
	"strconv"

	"github.com/go-redis/redis/v9"
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

func NewRedisClient(ctx context.Context) *redis.Client {
	redisDb, _ := strconv.Atoi(database)
	options := &redis.Options{
		Addr:     addr,
		Password: password,
		DB:       redisDb,
	}
	redisClient := redis.NewClient(options)
	return redisClient
}
