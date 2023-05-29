package events

import (
	"sync"
	"time"
)

type Event struct {
	Name      string    `json:"name"`
	Payload   any       `json:"payload"`
	CreatedAt time.Time `json:"createdAt"`
}

type Handler interface {
	Handle(event Event, wg *sync.WaitGroup)
}

type DispatcherInterface interface {
	Register(name string, handler Handler) error
	Dispatch(event Event) error
	Remove(name string, handler Handler) error
	Has(name string, handler Handler) bool
	Total(name string) int
	GetByIndex(name string, index int) Handler
	Clear()
}
