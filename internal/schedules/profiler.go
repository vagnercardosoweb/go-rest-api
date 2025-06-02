package schedules

import (
	"fmt"
	"runtime"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

const megaBytes = 1 << 20 // 1MB = 1024 * 1024 bytes

func runProfiler(e *Scheduler) error {
	if !env.GetAsBool("PROFILER_ENABLED") {
		return nil
	}

	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)

	e.logger.
		AddField("memoryUsed", fmt.Sprintf("%vmb", m.Alloc/megaBytes)).
		AddField("memoryAcquired", fmt.Sprintf("%vmb", m.Sys/megaBytes)).
		AddField("numGoroutine", runtime.NumGoroutine()).
		Info("PROFILER")

	return nil
}
