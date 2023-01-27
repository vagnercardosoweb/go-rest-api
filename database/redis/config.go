package redis

import (
	"fmt"
	"rest-api/shared"
	"strconv"

	libRedis "github.com/go-redis/redis/v9"
)

func NewConfig() *libRedis.Options {
	addr := fmt.Sprintf(
		"%s:%s",
		shared.EnvRequiredByName("REDIS_HOST"),
		shared.EnvRequiredByName("REDIS_PORT"),
	)

	database, _ := strconv.Atoi(shared.EnvRequiredByName("REDIS_DATABASE", "0"))
	password := shared.EnvRequiredByName("REDIS_PASSWORD")

	return &libRedis.Options{
		Addr:     addr,
		Password: password,
		DB:       database,
	}
}
