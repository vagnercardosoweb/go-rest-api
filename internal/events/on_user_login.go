package events

import (
	"github.com/vagnercardosoweb/go-rest-api/pkg/events"
)

type OnUserLoginEvent struct{ *Manager }

func NewOnUserLoginEvent(m *Manager) *OnUserLoginEvent {
	return &OnUserLoginEvent{m}
}

func (e *OnUserLoginEvent) Handle(event *events.Event) error {
	m := e.Manager.Clone(event.TraceId)
	input := event.Input.(OnUserLoginInput)

	_, err := m.pgClient.Exec(
		`
			UPDATE "users"
			SET
				"last_login_at" = NOW(),
				"last_login_agent" = $3,
				"last_login_ip" = $2
			WHERE
				"id" = $1;
		`,
		input.UserId,
		input.IpAddress,
		input.UserAgent,
	)

	return err
}

type OnUserLoginInput struct {
	UserId    string `json:"userId"`
	IpAddress string `json:"-"`
	UserAgent string `json:"-"`
	TraceId   string `json:"-"`
}

func (m *Manager) OnUserLogin(input OnUserLoginInput) {
	m.Dispatch(&events.Event{
		Name:    OnUserLoginName,
		TraceId: input.TraceId,
		Input:   input,
	})
}
