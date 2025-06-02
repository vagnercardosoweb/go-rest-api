package events

import (
	"time"
)

type Event struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	RequestId string    `json:"-"`
	Input     any       `json:"input"`
}

type Handler interface {
	Handle(event *Event) error
}

type DispatcherInterface interface {
	Register(name string, handler Handler) error
	Dispatch(event *Event) error
	Remove(name string, handler Handler) error
	Has(name string, handler Handler) bool
	Total(name string) int
	GetByIndex(name string, index int) Handler
	Clear()
}

type Dispatcher struct {
	handlers map[string][]Handler
}
