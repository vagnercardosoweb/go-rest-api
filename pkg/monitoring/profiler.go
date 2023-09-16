package monitoring

import (
	"fmt"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"runtime"
	"time"
)

func runProfiler(logger *logger.Logger) {
	m := &runtime.MemStats{}

	for {
		runtime.ReadMemStats(m)

		logger.
			WithID("MONITORING").
			AddMetadata("memory_used", m.Alloc).
			AddMetadata("memory_used_mb", fmt.Sprintf("%vmb", m.Alloc/config.MegaBytes)).
			AddMetadata("goroutine", runtime.NumGoroutine()).
			AddMetadata("memory_acquired_mb", fmt.Sprintf("%vmb", m.Sys/config.MegaBytes)).
			AddMetadata("memory_acquired", m.Sys).
			Info("PROFILER")

		time.Sleep(time.Second * 30)
	}
}

func RunProfiler(logger *logger.Logger) {
	go runProfiler(logger)
}
