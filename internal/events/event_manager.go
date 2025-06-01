package events

import (
	"github.com/vagnercardosoweb/go-rest-api/pkg/events"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
)

type EventManager struct {
	dispatcher  events.DispatcherInterface
	pgClient    *postgres.Client
	redisClient *redis.Client
	logger      *logger.Logger
}

func NewManager(pgClient *postgres.Client, redisClient *redis.Client) *EventManager {
	e := &EventManager{
		logger:      pgClient.GetLogger(),
		dispatcher:  events.NewDispatcher(),
		redisClient: redisClient,
		pgClient:    pgClient,
	}

	e.register(AfterLoginName, &AfterLogin{e})

	return e
}

func (e *EventManager) WithLogger(logger *logger.Logger) *EventManager {
	return NewManager(e.pgClient.WithLogger(logger), e.redisClient)
}

func (e *EventManager) dispatch(event *events.Event) {
	if err := e.dispatcher.Dispatch(event); err != nil {
		e.logger.
			AddMetadata("eventName", event.Name).
			AddMetadata("originalError", err).
			Error("EVENT_MANAGER_DISPATCH_ERROR")
	}
}

func (e *EventManager) register(name string, handler events.Handler) {
	if err := e.dispatcher.Register(name, handler); err != nil {
		e.logger.
			AddMetadata("eventName", name).
			AddMetadata("originalError", err).
			Error("EVENT_MANAGER_REGISTER_ERROR")
	}
}
