package events

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

const CtxKey = "EventManagerKey"

func FromGin(c *gin.Context) *EventManager {
	return c.MustGet(CtxKey).(*EventManager)
}

func FromCtx(c context.Context) *EventManager {
	value, exists := c.Value(CtxKey).(*EventManager)

	if !exists {
		panic(fmt.Errorf(`context key "%s" does not exist`, CtxKey))
	}

	return value
}
