package events

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

const CtxKey = "EventManagerKey"

func FromGin(c *gin.Context) *Manager {
	return c.MustGet(CtxKey).(*Manager)
}

func FromCtx(c context.Context) *Manager {
	value, exists := c.Value(CtxKey).(*Manager)

	if !exists {
		panic(fmt.Errorf(`context key "%s" does not exist`, CtxKey))
	}

	return value
}
