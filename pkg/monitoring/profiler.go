package monitoring

import (
	"fmt"
	"runtime"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

const megaBytes = 1 << 20

func runProfiler(logger *logger.Logger) {
	m := &runtime.MemStats{}

	for {
		runtime.ReadMemStats(m)

		logger.
			WithId("MONITORING").
			WithoutRedact().
			AddMetadata("memoryUsed", fmt.Sprintf("%vmb", m.Alloc/megaBytes)).
			AddMetadata("memoryAcquired", fmt.Sprintf("%vmb", m.Sys/megaBytes)).
			AddMetadata("numGoroutine", runtime.NumGoroutine()).
			Info("PROFILER")

		time.Sleep(time.Second * 30)
	}
}

func RunProfiler(logger *logger.Logger) {
	go runProfiler(logger)
}
