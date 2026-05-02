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
	fields     map[string]any
	mu         *sync.Mutex
}

type Output struct {
	Id          string         `json:"id"`
	Level       level          `json:"level"`
	Hostname    string         `json:"hostname"`
	Timestamp   time.Time      `json:"timestamp"`
	Environment string         `json:"environment"`
	Message     string         `json:"message"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}
