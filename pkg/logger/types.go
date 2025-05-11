package logger

import (
	"sync"
	"time"
)

type level string

const (
	LevelInfo  level = "INFO"
	LevelError level = "ERROR"
)

type Logger struct {
	id         string
	enabled    bool
	redactKeys []string
	metadata   map[string]any
	mu         *sync.Mutex
}

type Output struct {
	Id          string         `json:"id"`
	Level       level          `json:"level"`
	Timestamp   time.Time      `json:"timestamp"`
	Message     string         `json:"message"`
	Environment string         `json:"environment"`
	Hostname    string         `json:"hostname"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}
