package schedules

import (
	"sync"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/postgres"
	"github.com/vagnercardosoweb/go-rest-api/pkg/redis"
)

type Job func(s *Scheduler) error
type Scheduler struct {
	logger      *logger.Logger
	pgClient    *postgres.Client
	cacheClient *redis.Client
	wg          sync.WaitGroup
	sleep       time.Duration
	jobs        []Job
}
