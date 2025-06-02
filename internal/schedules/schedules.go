package schedules

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
	"github.com/vagnercardosoweb/go-rest-api/pkg/slack"
)

func New(
	pgClient *postgres.Client,
	redisClient *redis.Client,
) *Scheduler {
	s := &Scheduler{
		pgClient:    pgClient,
		cacheClient: redisClient,
		logger:      pgClient.Logger().WithId("SCHEDULER"),
		sleep:       env.GetSchedulerSleep(),
		wg:          sync.WaitGroup{},
		jobs:        make([]Job, 0),
	}

	s.AddJob(runProfiler)

	return s
}

func (s *Scheduler) AddJob(job Job) {
	s.jobs = append(s.jobs, job)
}

func (s *Scheduler) Run() {
	if len(s.jobs) == 0 {
		return
	}

	go func(s *Scheduler) {
		ticket := time.NewTicker(s.sleep)
		defer ticket.Stop()

		for range ticket.C {
			s.wg.Add(len(s.jobs))

			for _, job := range s.jobs {
				go func(job Job) {
					defer s.recover()
					defer s.wg.Done()

					err := job(s)
					s.notifyError(err, false)
				}(job)
			}

			s.wg.Wait()
		}
	}(s)
}

func (s *Scheduler) notifyError(err any, isPanic bool) {
	if err == nil {
		return
	}

	_, file, line, _ := runtime.Caller(3)
	caller := fmt.Sprintf("%s:%d", file, line)

	var message string
	traceId := uuid.New().String()

	if isPanic {
		message = "A panic error was received when executing job processing"
	} else {
		message = "An error was received when executing job processing"
	}

	s.logger.
		AddField("traceId", traceId).
		AddField("error", err).
		Error(message)

	go func() {
		_ = slack.NewAlert().
			AddField("caller", caller, false).
			AddField("traceId", traceId, false).
			AddField("message", message, false).
			AddError("error", err).
			Send()
	}()
}

func (s *Scheduler) recover() {
	s.notifyError(recover(), true)
}
