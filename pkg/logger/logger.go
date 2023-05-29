package logger

import (
	"encoding/json"
	"fmt"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"log"
	"os"
	"sync"
	"time"
)

var logger = log.New(os.Stdout, "", 0)

type level string

const (
	LevelInfo     level = "INFO"
	LevelDebug    level = "DEBUG"
	LevelWarn     level = "WARN"
	LevelCritical level = "CRITICAL"
	LevelError    level = "ERROR"
)

type Logger struct {
	id       string
	metadata map[string]any
	mu       *sync.Mutex
}

type output struct {
	Id        string         `json:"id"`
	Level     level          `json:"level"`
	Message   string         `json:"message"`
	Pid       int            `json:"pid"`
	Hostname  string         `json:"hostname"`
	Timestamp time.Time      `json:"timestamp"`
	Metadata  map[string]any `json:"metadata,omitempty"`
}

func New() *Logger {
	return &Logger{
		id:       "APP",
		metadata: make(map[string]any),
		mu:       new(sync.Mutex),
	}
}

func (*Logger) WithID(id string) *Logger {
	l := New()
	l.id = id
	return l
}

func (l *Logger) WithMetadata(metadata map[string]any) *Logger {
	nl := New()
	nl.metadata = metadata
	nl.id = l.id
	return nl
}

func (l *Logger) AddMetadata(key string, value any) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.metadata[key] = value
	return l
}

func (l *Logger) Info(message string, arguments ...any) {
	l.Log(LevelInfo, message, arguments...)
}

func (l *Logger) Warn(message string, arguments ...any) {
	l.Log(LevelWarn, message, arguments...)
}

func (l *Logger) Debug(message string, arguments ...any) {
	l.Log(LevelDebug, message, arguments...)
}

func (l *Logger) Critical(message string, arguments ...any) {
	l.Log(LevelCritical, message, arguments...)
}

func (l *Logger) Error(message string, arguments ...any) {
	l.Log(LevelError, message, arguments...)
}

func (l *Logger) Log(level level, message string, arguments ...any) {
	if len(arguments) > 0 {
		message = fmt.Sprintf(message, arguments...)
	}
	logAsJson, _ := json.Marshal(output{
		Id:        l.id,
		Level:     level,
		Message:   message,
		Timestamp: time.Now().UTC(),
		Metadata:  l.metadata,
		Pid:       config.Pid,
		Hostname:  config.Hostname,
	})
	l.metadata = make(map[string]any)
	logger.Println(string(logAsJson))
}
