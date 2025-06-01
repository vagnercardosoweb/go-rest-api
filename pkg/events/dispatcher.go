package events

import (
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

func NewDispatcher() DispatcherInterface {
	return &Dispatcher{handlers: make(map[string][]Handler)}
}

func (d *Dispatcher) Total(name string) int {
	if name == "" {
		return len(d.handlers)
	}

	return len(d.handlers[name])
}

func (d *Dispatcher) GetByIndex(name string, index int) Handler {
	return d.handlers[name][index]
}

func (d *Dispatcher) Register(name string, handler Handler) error {
	if d.Has(name, handler) {
		return errors.FromMessage(`event handler "%s" already registered`, name)
	}

	d.handlers[name] = append(d.handlers[name], handler)
	return nil
}

func (d *Dispatcher) Dispatch(event *Event) error {
	if handlers, ok := d.handlers[event.Name]; ok {
		wg := &sync.WaitGroup{}
		wg.Add(len(handlers))

		for i, handler := range handlers {
			go func() {
				defer wg.Done()

				if event.CreatedAt.IsZero() {
					event.CreatedAt = time.Now()
				}

				if event.CorrelationId == "" {
					event.CorrelationId = fmt.Sprintf(
						"%s_HANDLER_%d",
						strings.ToUpper(event.Name),
						i+1,
					)
				}

				handler.Handle(event)
			}()
		}

		wg.Wait()
	}

	return nil
}

func (d *Dispatcher) Remove(name string, handler Handler) error {
	if _, ok := d.handlers[name]; ok {
		index := slices.Index(d.handlers[name], handler)

		if index != -1 {
			d.handlers[name] = slices.Delete(d.handlers[name], index, index+1)
			return nil
		}
	}

	return nil
}

func (d *Dispatcher) Has(name string, handler Handler) bool {
	if _, ok := d.handlers[name]; ok {
		if slices.Contains(d.handlers[name], handler) {
			return true
		}
	}

	return false
}

func (d *Dispatcher) Clear() {
	d.handlers = make(map[string][]Handler)
}
