package redis

import (
	"fmt"
	"strconv"

	libRedis "github.com/go-redis/redis/v9"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

func newConfig() *libRedis.Options {
	addr := fmt.Sprintf(
		"%s:%s",
		env.Required("REDIS_HOST"),
		env.Required("REDIS_PORT"),
	)

	database, _ := strconv.Atoi(env.Required("REDIS_DATABASE"))
	password := env.Required("REDIS_PASSWORD")

	return &libRedis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	}
}
