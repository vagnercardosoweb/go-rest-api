package events

import (
	eventspkg "github.com/vagnercardosoweb/go-rest-api/pkg/events"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
)

type EventManager struct {
	pgClient    *postgres.Client
	redisClient *redis.Client
	dispatcher  eventspkg.DispatcherInterface
	logger      *logger.Logger
}

func New(pgClient *postgres.Client, redisClient *redis.Client) *EventManager {
	event := &EventManager{
		pgClient:    pgClient,
		redisClient: redisClient,
		dispatcher:  eventspkg.NewDispatcher(),
		logger:      pgClient.GetLogger(),
	}

	RegisterAfterLogin(event)

	return event
}

func (e *EventManager) Dispatch(event *eventspkg.Event) {
	err := e.dispatcher.Dispatch(event)
	if err != nil {
		e.logger.
			WithoutRedact().
			AddMetadata("originalError", err).
			AddMetadata("eventName", event.Name).
			Error("failed to dispatch event")

	}
}

func (e *EventManager) Register(name string, handler eventspkg.Handler) {
	err := e.dispatcher.Register(name, handler)
	if err != nil {
		e.logger.
			WithoutRedact().
			AddMetadata("originalError", err).
			AddMetadata("eventName", name).
			Error("failed to register event")
	}
}
