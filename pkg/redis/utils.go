package redis

import (
	"context"
	"fmt"
)

const CtxKey = "RedisClientKey"

func FromCtx(c context.Context) *Client {
	value, exists := c.Value(CtxKey).(*Client)

	if !exists {
		panic(fmt.Errorf(`context key "%s" does not exist`, CtxKey))
	}

	return value
}
