package logger

import (
	"sync"
	"time"
)

type level string

const (
	LevelInfo     level = "INFO"
	LevelDebug    level = "DEBUG"
	LevelWarn     level = "WARN"
	LevelCritical level = "CRITICAL"
	LevelError    level = "ERROR"
)

type Logger struct {
	id              string
	metadata        map[string]any
	isLoggerEnabled bool
	isDebugEnabled  bool
	redactKeys      []string
	skipRedact      bool
	mu              *sync.Mutex
}

type Output struct {
	Id          string         `json:"id"`
	Level       level          `json:"level"`
	Pid         int            `json:"pid"`
	Hostname    string         `json:"hostname"`
	Timestamp   time.Time      `json:"timestamp"`
	Environment string         `json:"environment"`
	Message     string         `json:"message"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}
