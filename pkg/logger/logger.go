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
		id:         "APP",
		metadata:   make(map[string]any),
		redactKeys: strings.Split(env.GetAsString("REDACT_KEYS", ""), ","),
		enabled:    env.GetAsBool("LOGGER_ENABLED", "true"),
		mu:         new(sync.Mutex),
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
	return l.AddMetadata("skipRedact", "true")
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

func (l *Logger) Error(message string, arguments ...any) {
	l.Log(LevelError, message, arguments...)
}

func (l *Logger) Log(level level, message string, arguments ...any) {
	if !l.enabled {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if len(arguments) > 0 {
		message = fmt.Sprintf(message, arguments...)
	}

	if l.isRedact() {
		l.metadata = utils.RedactKeys(l.metadata, l.redactKeys)
	}

	logAsJson, _ := json.Marshal(Output{
		Id:          l.id,
		Level:       level,
		Hostname:    utils.Hostname,
		Environment: env.GetAppEnv(),
		Timestamp:   time.Now().UTC(),
		Metadata:    l.metadata,
		Message:     message,
	})

	l.metadata = make(map[string]any)
	logger.Println(string(logAsJson))
}

func (l *Logger) isRedact() bool {
	if _, ok := l.metadata["skipRedact"]; ok {
		delete(l.metadata, "skipRedact")

		return false
	}

	return len(l.metadata) > 0 &&
		len(l.redactKeys) > 0
}

const CtxKey = "LoggerKey"

func GetFromCtxOrPanic(ctx context.Context) *Logger {
	l, ok := ctx.Value(CtxKey).(*Logger)

	if !ok {
		panic(fmt.Sprintf(`context key "%s" does not exist`, CtxKey))
	}

	return l
}
