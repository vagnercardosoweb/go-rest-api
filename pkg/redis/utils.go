package redis

import "context"

const CtxKey = "RedisClientCtxKey"

func GetFromCtxOrPanic(c context.Context) *Client {
	value, exists := c.Value(CtxKey).(*Client)
	if !exists {
		panic("redis client not found in context")
	}
	return value
}
