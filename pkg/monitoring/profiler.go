package monitoring

import (
	"fmt"
	"runtime"
	"time"

	logger "github.com/vagnercardosoweb/go-rest-api/pkg/logger_new"
)

const megabyte = 1024 * 1024

func runProfiler() {
	m := &runtime.MemStats{}

	for {
		runtime.ReadMemStats(m)

		logger.New().
			WithID("MONITORING").
			WithMetadata("memory_used", m.Alloc).
			WithMetadata("memory_used_mb", fmt.Sprintf("%vmb", m.Alloc/megabyte)).
			WithMetadata("goroutine", runtime.NumGoroutine()).
			WithMetadata("memory_acquired_mb", fmt.Sprintf("%vmb", m.Sys/megabyte)).
			WithMetadata("memory_acquired", m.Sys).
			Info("Run profiler")

		time.Sleep(time.Second * 30)
	}
}

func RunProfiler() {
	go runProfiler()
}
