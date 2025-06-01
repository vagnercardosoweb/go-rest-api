package schedules

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/slack"
)

type Job func(s *Scheduler) error

type Scheduler struct {
	logger      *logger.Logger
	pgClient    *postgres.Client
	cacheClient *redis.Client
	waiter      sync.WaitGroup
	sleep       time.Duration
	jobs        []Job
}

func New(
	pgClient *postgres.Client,
	redisClient *redis.Client,
	logger *logger.Logger,
) *Scheduler {
	return &Scheduler{
		pgClient:    pgClient,
		cacheClient: redisClient,
		logger:      logger.WithId("SCHEDULER"),
		sleep:       env.GetSchedulerSleep(),
		waiter:      sync.WaitGroup{},
		jobs:        make([]Job, 0),
	}
}

func (s *Scheduler) AddJob(job Job) {
	s.jobs = append(s.jobs, job)
}

func (s *Scheduler) Run() {
	if len(s.jobs) == 0 {
		return
	}

	for {
		time.Sleep(s.sleep * time.Second)
		s.waiter.Add(len(s.jobs))

		for _, job := range s.jobs {
			go func(job Job) {
				defer s.recover()
				defer s.waiter.Done()

				if err := job(s); err != nil {
					s.notifySlackOfError(err, false)
				}
			}(job)
		}

		s.waiter.Wait()
	}
}

func (s *Scheduler) notifySlackOfError(err any, isPanic bool) {
	_, file, line, _ := runtime.Caller(3)
	caller := fmt.Sprintf("%s:%d", file, line)

	var trackId string
	if v, ok := err.(*errors.Input); ok {
		trackId = v.RequestId
	} else {
		trackId = uuid.NewString()
	}

	var message string
	if isPanic {
		message = "A panic error was received when executing job processing"
	} else {
		message = "An error was received when executing job processing"
	}

	s.logger.
		AddMetadata("trackId", trackId).
		AddMetadata("error", err).
		Error(message)

	go func() {
		_ = slack.NewAlert().
			AddField("caller", caller, false).
			AddField("trackId", trackId, false).
			AddField("message", message, false).
			AddError("error", err).
			Send()
	}()
}

func (s *Scheduler) recover() {
	if r := recover(); r != nil {
		s.notifySlackOfError(r, true)
	}
}
