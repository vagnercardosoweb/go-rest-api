package events

import (
	"fmt"
	"sync"
)

type Dispatcher struct {
	handlers map[string][]Handler
}

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
		return fmt.Errorf("event handler [%s] already registered", name)
	}
	d.handlers[name] = append(d.handlers[name], handler)
	return nil
}

func (d *Dispatcher) Dispatch(event Event) error {
	if handlers, ok := d.handlers[event.Name]; ok {
		wg := &sync.WaitGroup{}
		for _, handler := range handlers {
			wg.Add(1)
			go handler.Handle(event, wg)
		}
		wg.Wait()
	}
	return nil
}

func (d *Dispatcher) Remove(name string, handler Handler) error {
	if _, ok := d.handlers[name]; ok {
		for i, h := range d.handlers[name] {
			if h == handler {
				d.handlers[name] = append(d.handlers[name][:i], d.handlers[name][i+1:]...)
				return nil
			}
		}
	}
	return nil
}

func (d *Dispatcher) Has(name string, handler Handler) bool {
	if _, ok := d.handlers[name]; ok {
		for _, h := range d.handlers[name] {
			if h == handler {
				return true
			}
		}
	}
	return false
}

func (d *Dispatcher) Clear() {
	d.handlers = make(map[string][]Handler)
}
