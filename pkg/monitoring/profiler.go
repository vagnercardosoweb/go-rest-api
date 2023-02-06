package monitoring

import (
	"fmt"
	"runtime"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

const megabyte = 1024 * 1024

func runProfiler() {
	m := &runtime.MemStats{}

	for {
		runtime.ReadMemStats(m)

		logger.Log(logger.Input{
			Id:      "MONITORING",
			Level:   logger.DEBUG,
			Message: "Run profiler",
			Metadata: logger.Metadata{
				"memory_used":        m.Alloc,
				"memory_used_mb":     fmt.Sprintf("%v mb", m.Alloc/megabyte),
				"goroutine":          runtime.NumGoroutine(),
				"memory_acquired_mb": fmt.Sprintf("%v mb", m.Sys/megabyte),
				"memory_acquired":    m.Sys,
			},
		})

		time.Sleep(time.Second * 30)
	}
}

func RunProfiler() {
	go runProfiler()
}
