package events

import (
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/events"
)

type AfterLogin struct{ *EventManager }

func (e *AfterLogin) Handle(event *events.Event) {
	e.logger.Info("Handling after login event")
}

type AfterLoginInput struct {
	UserId string `json:"user_id"`
}

func (e *EventManager) AfterLogin(input AfterLoginInput) *events.Event {
	event := &events.Event{
		Name:          AfterLoginName,
		CorrelationId: e.logger.GetId(),
		CreatedAt:     time.Now(),
		Input:         input,
	}

	e.dispatch(event)
	return event
}
