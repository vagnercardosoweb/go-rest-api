package shared

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"
)

type Logger struct {
	Id       string                 `json:"id"`
	Level    string                 `json:"level"`
	Message  string                 `json:"message"`
	Time     time.Time              `json:"time"`
	Pid      int                    `json:"pid"`
	Hostname string                 `json:"hostname"`
	Metadata map[string]interface{} `json:"metadata"`
}

var pid = os.Getpid()
var hostname, _ = os.Hostname()

func GetLogger() *Logger {
	logger := NewLogger(Logger{Id: "APP"})
	return logger
}

func NewLogger(input Logger) *Logger {
	if input.Id == "" {
		input.Id = "APP"
	}
	if input.Pid == 0 {
		input.Pid = pid
	}
	if input.Hostname == "" {
		input.Hostname = hostname
	}
	if input.Level == "" {
		input.Level = "INFO"
	}
	return &input
}

func (l *Logger) Log(level string, message string) {
	l.Time = time.Now().UTC()
	l.Level = strings.ToUpper(level)
	l.Message = message
	structToJson, _ := json.Marshal(l)
	log.New(os.Stdout, "", 0).Printf(string(structToJson))
	l.Metadata = nil
}

func (l *Logger) AddMetadata(name string, value interface{}) *Logger {
	if l.Metadata == nil {
		l.Metadata = make(map[string]interface{})
	}
	l.Metadata[name] = value
	return l
}

func (l *Logger) Info(message string) {
	l.Log("INFO", message)
}

func (l *Logger) Warning(message string) {
	l.Log("WARNING", message)
}

func (l *Logger) Error(message string) {
	l.Log("ERROR", message)
}

func (l *Logger) Critical(message string) {
	l.Log("CRITICAL", message)
}

func (l *Logger) Debug(message string) {
	l.Log("DEBUG", message)
}
