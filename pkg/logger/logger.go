package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"github.com/vagnercardosoweb/go-rest-api/pkg/utils"
)

var logger = log.New(os.Stdout, "", 0)

func New() *Logger {
	return &Logger{
		id:              "APP",
		metadata:        make(map[string]any),
		isDebugEnabled:  env.GetAsBool("LOGGER_DEBUG", "true"),
		isLoggerEnabled: env.GetAsBool("LOGGER_ENABLED", "true"),
		redactKeys:      strings.Split(env.GetAsString("REDACT_KEYS", ""), ","),
		mu:              new(sync.Mutex),
	}
}

func (*Logger) WithId(id string) *Logger {
	l := New()
	l.id = id
	return l
}

func (l *Logger) GetId() string {
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
	if !l.isDebugEnabled {
		return
	}

	l.Log(LevelDebug, message, arguments...)
}

func (l *Logger) Critical(message string, arguments ...any) {
	l.Log(LevelCritical, message, arguments...)
}

func (l *Logger) Error(message string, arguments ...any) {
	l.Log(LevelError, message, arguments...)
}

func (l *Logger) Log(level level, message string, arguments ...any) {
	if !l.isLoggerEnabled {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if len(arguments) > 0 {
		message = fmt.Sprintf(message, arguments...)
	}

	if len(l.metadata) > 0 && len(l.redactKeys) > 0 && l.skipRedact == false {
		startRedact := time.Now()
		l.metadata = utils.RedactKeys(l.metadata, l.redactKeys)
		l.metadata["redactTime"] = time.Since(startRedact).String()
	}

	logAsJson, _ := json.Marshal(Output{
		Id:          l.id,
		Level:       level,
		Environment: env.GetAppEnv(),
		Pid:         utils.Pid,
		Hostname:    utils.Hostname,
		Timestamp:   time.Now().UTC(),
		Message:     message,
		Metadata:    l.metadata,
	})

	l.skipRedact = false
	l.metadata = make(map[string]any)

	logger.Println(string(logAsJson))
}

const CtxKey = "LoggerCtxKey"

func GetFromCtxOrPanic(ctx context.Context) *Logger {
	l, ok := ctx.Value(CtxKey).(*Logger)
	if !ok {
		panic("Logger not found in context")
	}
	return l
}
