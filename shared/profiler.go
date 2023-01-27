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

		logger.
			AddMetadata("NumGoroutine", runtime.NumGoroutine()).
			AddMetadata("MemoryUsed", m.Alloc).
			AddMetadata("MemoryUsedInMb", fmt.Sprintf("%v mb", m.Alloc/MEGABYTE)).
			AddMetadata("MemoryAcquired", m.Sys).
			AddMetadata("MemoryAcquiredInMb", fmt.Sprintf("%v mb", m.Sys/MEGABYTE)).
			Debug("Running profiler")

		time.Sleep(time.Second * 30)
	}
}

func StartProfiler() {
	go runProfiler()
}
