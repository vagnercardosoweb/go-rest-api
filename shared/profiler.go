package shared

import (
	"fmt"
	"runtime"
	"time"
)

const MEGABYTE = 1024 * 1024

var logger = NewLogger(Logger{Id: "PROFILE"})

func runProfiler() {
	m := &runtime.MemStats{}

	for {
		runtime.ReadMemStats(m)

		logger.AddMetadata("NumGoroutine", runtime.NumGoroutine())
		logger.AddMetadata("MemoryUsed", m.Alloc)
		logger.AddMetadata("MemoryUsedInMb", fmt.Sprintf("%v mb", m.Alloc/MEGABYTE))
		logger.AddMetadata("MemoryAcquired", m.Sys)
		logger.AddMetadata("MemoryAcquiredInMb", fmt.Sprintf("%v mb", m.Sys/MEGABYTE))
		logger.Debug("Running profiler")

		time.Sleep(time.Second * 30)
	}
}

func StartProfiler() {
	go runProfiler()
}
