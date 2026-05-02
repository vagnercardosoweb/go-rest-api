package events

import (
	"github.com/vagnercardosoweb/go-rest-api/internal/repositories/user"
	"github.com/vagnercardosoweb/go-rest-api/pkg/events"
)

type OnUserLoginEvent struct{ *Manager }

func NewOnUserLoginEvent(m *Manager) *OnUserLoginEvent {
	return &OnUserLoginEvent{m}
}

func (e *OnUserLoginEvent) Handle(event *events.Event) error {
	m := e.Manager.Clone(event.TraceId)

	repo := user.New(m.pgClient)
	input := event.Input.(OnUserLoginInput)

	err := repo.UpdateLastLogin(&user.UpdateLastLoginInput{
		UserId:    input.UserId,
		IpAddress: input.IpAddress,
		UserAgent: input.UserAgent,
	})

	return err
}

type OnUserLoginInput struct {
	UserId    string
	IpAddress string
	UserAgent string
	TraceId   string
}

func (m *Manager) OnUserLogin(input OnUserLoginInput) {
	m.Dispatch(&events.Event{
		Name:    OnUserLoginName,
		TraceId: input.TraceId,
		Input:   input,
	})
}
