package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

var (
	pid      int
	hostname string
	logger   *log.Logger
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
	id         string
	metadata   map[string]any
	redactKeys []string
	skipRedact bool
	mu         *sync.Mutex
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

func init() {
	logger = log.New(os.Stdout, "", 0)
	hostname, _ = os.Hostname()
	pid = os.Getpid()
}

func New() *Logger {
	return &Logger{
		id:         "APP",
		metadata:   make(map[string]any),
		redactKeys: strings.Split(env.Get("OBFUSCATE_KEYS", ""), ","),
		mu:         new(sync.Mutex),
	}
}

func (*Logger) WithID(id string) *Logger {
	l := New()
	l.id = id
	return l
}

func (l *Logger) GetID() string {
	return l.id
}

func (l *Logger) WithoutRedact() *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.skipRedact = true
	return l
}

func (l *Logger) WithMetadata(metadata map[string]any) *Logger {
	for key, value := range metadata {
		l.AddMetadata(key, value)
	}
	return l
}

func (l *Logger) AddMetadata(key string, value any) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := value.(*errors.Input); !ok {
		if err, ok := value.(error); ok {
			value = err.Error()
		}
	}
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
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(arguments) > 0 {
		message = fmt.Sprintf(message, arguments...)
	}
	if len(l.metadata) > 0 && len(l.redactKeys) > 0 && l.skipRedact == false {
		startRedact := time.Now()
		l.metadata = redactKeys(l.metadata, l.redactKeys)
		elapsedRedact := time.Since(startRedact)
		l.metadata["redactTime"] = elapsedRedact.String()
	}
	logAsJson, _ := json.Marshal(Output{
		Id:          l.id,
		Level:       level,
		Environment: env.Get("APP_ENV", "local"),
		Pid:         pid,
		Hostname:    hostname,
		Timestamp:   time.Now().UTC(),
		Message:     message,
		Metadata:    l.metadata,
	})
	l.skipRedact = false
	l.metadata = make(map[string]any)
	logger.Println(string(logAsJson))
}
