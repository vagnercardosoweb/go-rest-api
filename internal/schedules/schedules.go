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
	pgClient    *postgres.Client
	cacheClient *redis.Client
	logger      *logger.Logger
	sleep       time.Duration
	waiter      sync.WaitGroup
	jobs        []Job
}

func New(
	pgClient *postgres.Client,
	redisClient *redis.Client,
	logger *logger.Logger,
) *Scheduler {
	scheduler := &Scheduler{
		pgClient:    pgClient,
		cacheClient: redisClient,
		waiter:      sync.WaitGroup{},
		sleep:       env.GetSchedulerSleep(),
		logger:      logger,
	}
	return scheduler
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
				defer s.waiter.Done()
				defer s.recover()

				if err := job(s); err != nil {
					s.sendErrorToSlack(err, false)
				}
			}(job)
		}

		s.waiter.Wait()
	}
}

func (s *Scheduler) sendErrorToSlack(err any, panic bool) {
	_, file, line, _ := runtime.Caller(3)
	caller := fmt.Sprintf("%s:%d", file, line)

	var trackId string
	if v, ok := err.(*errors.Input); ok {
		trackId = v.ErrorId
	} else {
		trackId = uuid.NewString()
	}

	var message string
	if panic {
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
			AddField("Caller", caller, false).
			AddField("TrackId", trackId, false).
			AddField("Message", message, false).
			AddError("Error", err).
			Send()
	}()
}

func (s *Scheduler) recover() {
	if r := recover(); r != nil {
		s.sendErrorToSlack(r, true)
	}
}

func (s *Scheduler) AddJob(job Job) {
	s.jobs = append(s.jobs, job)
}
