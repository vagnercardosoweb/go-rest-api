package monitoring

import (
	"fmt"
	"runtime"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

const MEGABYTES = 1024 * 1024

func runProfiler(logger *logger.Logger) {
	m := &runtime.MemStats{}

	for {
		runtime.ReadMemStats(m)

		logger.
			WithID("MONITORING").
			AddMetadata("memory_used", m.Alloc).
			AddMetadata("memory_used_mb", fmt.Sprintf("%vmb", m.Alloc/MEGABYTES)).
			AddMetadata("goroutine", runtime.NumGoroutine()).
			AddMetadata("memory_acquired_mb", fmt.Sprintf("%vmb", m.Sys/MEGABYTES)).
			AddMetadata("memory_acquired", m.Sys).
			Info("PROFILER")

		time.Sleep(time.Second * 30)
	}
}

func RunProfiler(logger *logger.Logger) {
	go runProfiler(logger)
}
