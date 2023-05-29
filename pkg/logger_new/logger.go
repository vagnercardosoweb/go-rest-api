package logger_new

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

type Logger struct {
	id       string
	metadata map[string]any
	mu       *sync.Mutex
}

type output struct {
	Id        string    `json:"id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Pid       int       `json:"pid"`
	Hostname  string    `json:"hostname"`
	Timestamp time.Time `json:"timestamp"`
	Metadata  any       `json:"metadata,omitempty"`
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

func (l *Logger) WithMetadata(key string, value any) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.metadata == nil {
		l.metadata = make(map[string]any)
	}
	l.metadata[key] = value
	return l
}

func (l *Logger) Info(message string, arguments ...any) {
	l.log("INFO", message, arguments...)
}

func (l *Logger) Warn(message string, arguments ...any) {
	l.log("WARN", message, arguments...)
}

func (l *Logger) Debug(message string, arguments ...any) {
	l.log("DEBUG", message, arguments...)
}

func (l *Logger) Error(message string, arguments ...any) {
	l.log("ERROR", message, arguments...)
}

func (l *Logger) log(level string, message string, arguments ...any) {
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
