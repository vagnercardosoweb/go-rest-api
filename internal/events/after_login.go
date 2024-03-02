package events

import (
	eventspkg "github.com/vagnercardosoweb/go-rest-api/pkg/events"
	"sync"
	"time"
)

type AfterLogin struct {
	eventManager *EventManager
}

const eventName = "afterLogin"

func RegisterAfterLogin(eventManager *EventManager) *AfterLogin {
	event := &AfterLogin{eventManager: eventManager}
	eventManager.Register(eventName, event)
	return event
}

func MakeAfterLogin(userId string) *eventspkg.Event {
	return &eventspkg.Event{
		Name:      eventName,
		CreatedAt: time.Now(),
		Payload:   userId,
	}
}

func (e *AfterLogin) Handle(event *eventspkg.Event, wg *sync.WaitGroup) {
	defer wg.Done()
}
