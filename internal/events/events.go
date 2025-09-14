package events

import (
	"github.com/vagnercardosoweb/go-rest-api/pkg/events"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/slack"
)

type Manager struct {
	pgClient    *postgres.Client
	dispatcher  events.DispatcherInterface
	redisClient *redis.Client
	logger      *logger.Logger
}

func NewManager(pgClient *postgres.Client, redisClient *redis.Client) *Manager {
	m := &Manager{
		logger:      pgClient.Logger(),
		dispatcher:  events.NewDispatcher(),
		redisClient: redisClient,
		pgClient:    pgClient,
	}

	m.Register(OnUserLoginName, NewOnUserLoginEvent(m))

	return m
}

func (m *Manager) Clone(traceId string) *Manager {
	if m.logger.GetId() == traceId {
		return m
	}

	l := m.logger.WithId(traceId)

	return &Manager{
		dispatcher:  m.dispatcher,
		pgClient:    m.pgClient.WithLogger(l),
		redisClient: m.redisClient,
		logger:      l,
	}
}

func (m *Manager) Dispatch(event *events.Event) {
	l := m.logger.WithId(event.TraceId)

	l.
		WithStruct(event).
		Info("EVENT_MANAGER_DISPATCH_EVENT")

	if err := m.dispatcher.Dispatch(event); err != nil {
		l.
			AddField("error", err).
			Error("EVENT_MANAGER_DISPATCH_ERROR")

		go func() {
			_ = slack.NewAlert().
				WithColor(slack.ColorError).
				AddField("eventName", event.Name, false).
				AddField("traceId", event.TraceId, false).
				AddField("message", err.Error(), false).
				Send()
		}()
	}
}

func (m *Manager) Register(name string, handler events.Handler) {
	if err := m.dispatcher.Register(name, handler); err != nil {
		m.logger.
			AddField("name", name).
			AddField("error", err).
			Error("EVENT_MANAGER_REGISTER_ERROR")
	}
}
