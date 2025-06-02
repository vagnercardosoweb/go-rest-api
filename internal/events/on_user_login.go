package events

import (
	"github.com/vagnercardosoweb/go-rest-api/pkg/events"
)

type OnUserLoginEvent struct{ *Manager }

func NewOnUserLoginEvent(m *Manager) *OnUserLoginEvent {
	return &OnUserLoginEvent{m}
}

func (e *OnUserLoginEvent) Handle(event *events.Event) error {
	m := e.Manager.Clone(event.RequestId)
	input := event.Input.(OnUserLoginInput)

	_, err := m.pgClient.Exec(updateLastLoginQuery, input.UserId, input.IpAddress, input.UserAgent)
	return err
}

var updateLastLoginQuery = `
	UPDATE "users"
	SET
		"last_login_at" = NOW(),
		"last_login_agent" = $3,
		"last_login_ip" = $2
	WHERE
		"id" = $1;
`

type OnUserLoginInput struct {
	UserId    string `json:"userId"`
	IpAddress string `json:"-"`
	UserAgent string `json:"-"`
	RequestId string `json:"-"`
}

func (m *Manager) SendOnUserLogin(input OnUserLoginInput) {
	event := &events.Event{
		Name:      OnUserLoginName,
		RequestId: input.RequestId,
		Input:     input,
	}

	m.Dispatch(event)
}
