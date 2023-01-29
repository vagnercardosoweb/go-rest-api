package monitoring

import (
	"fmt"
	"runtime"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

const megabyte = 1024 * 1024

var log = logger.New(logger.Input{Id: "MONITORING"})

func runProfiler() {
	m := &runtime.MemStats{}

	for {
		runtime.ReadMemStats(m)

		log.
			AddMetadata("NumGoroutine", runtime.NumGoroutine()).
			AddMetadata("MemoryUsed", m.Alloc).
			AddMetadata("MemoryUsedInMb", fmt.Sprintf("%v mb", m.Alloc/megabyte)).
			AddMetadata("MemoryAcquired", m.Sys).
			AddMetadata("MemoryAcquiredInMb", fmt.Sprintf("%v mb", m.Sys/megabyte)).
			Debug("Running profiler")

		time.Sleep(time.Second * 30)
	}
}

func RunProfiler() {
	go runProfiler()
}
