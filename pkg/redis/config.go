package redis

import (
	"fmt"
	"strconv"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"

	redisLib "github.com/go-redis/redis/v9"
)

func newConfig() *redisLib.Options {
	addr := fmt.Sprintf(
		"%s:%s",
		env.Required("REDIS_HOST"),
		env.Required("REDIS_PORT"),
	)

	database, _ := strconv.Atoi(env.Required("REDIS_DATABASE"))
	password := env.Required("REDIS_PASSWORD")

	return &redisLib.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	}
}
