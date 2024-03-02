package events

import (
	"context"
	"github.com/gin-gonic/gin"
)

const CtxKey = "EventManagerCtxKey"

func GetFromGinCtx(c *gin.Context) *EventManager {
	return c.MustGet(CtxKey).(*EventManager)
}

func GetFromCtxOrPanic(c context.Context) *EventManager {
	value, exists := c.Value(CtxKey).(*EventManager)
	if !exists {
		panic("event manager not found in context")
	}
	return value
}
