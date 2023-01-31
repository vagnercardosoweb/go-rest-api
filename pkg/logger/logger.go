package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
)

type (
	Metadata map[string]any
	Input    struct {
		Id        string
		Level     string
		Message   string
		Metadata  Metadata
		Arguments []any
	}
)

var logger = log.New(os.Stdout, "", 0)
var (
	LevelInfo     = "INFO"
	LevelWarn     = "WARN"
	LevelCritical = "CRITICAL"
	LevelError    = "ERROR"
	LevelDebug    = "DEBUG"
)

func Log(input Input) {
	if input.Id == "" {
		input.Id = "APP"
	}

	if input.Level == "" {
		input.Level = LevelInfo
	}

	logJson, _ := json.Marshal(struct {
		Id       string    `json:"id"`
		Level    string    `json:"level"`
		Message  string    `json:"message"`
		Pid      int       `json:"pid"`
		Hostname string    `json:"hostname"`
		Time     time.Time `json:"time"`
		Metadata Metadata  `json:"metadata"`
	}{
		Id:       input.Id,
		Level:    input.Level,
		Message:  fmt.Sprintf(input.Message, input.Arguments...),
		Time:     time.Now().UTC(),
		Metadata: input.Metadata,
		Pid:      config.Pid,
		Hostname: config.Hostname,
	})

	logger.Print(string(logJson))
}

func Info(message string, arguments ...any) {
	Log(Input{
		Level:     LevelInfo,
		Arguments: arguments,
		Message:   message,
	})
}

func Warn(message string, arguments ...any) {
	Log(Input{
		Level:     LevelWarn,
		Arguments: arguments,
		Message:   message,
	})
}

func Error(message string, arguments ...any) {
	Log(Input{
		Level:     LevelError,
		Arguments: arguments,
		Message:   message,
	})
}

func Critical(message string, arguments ...any) {
	Log(Input{
		Level:     LevelCritical,
		Arguments: arguments,
		Message:   message,
	})
}

func Debug(message string, arguments ...any) {
	Log(Input{
		Level:     LevelInfo,
		Arguments: arguments,
		Message:   message,
	})
}
