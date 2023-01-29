package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Input struct {
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
var logger = log.New(os.Stdout, "", 0)

func Get() *Input {
	return New(Input{Id: "APP"})
}

func New(input Input) *Input {
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

func (l *Input) Log(level string, message string, arguments ...interface{}) {
	l.Time = time.Now().UTC()
	l.Level = strings.ToUpper(level)
	l.Message = fmt.Sprintf(message, arguments...)
	structToJson, _ := json.Marshal(l)
	logger.Print(string(structToJson))
	l.Metadata = nil
}

func (l *Input) AddMetadata(name string, value interface{}) *Input {
	if l.Metadata == nil {
		l.Metadata = make(map[string]interface{})
	}
	l.Metadata[name] = value
	return l
}

func (l *Input) Info(message string, arguments ...interface{}) {
	l.Log("INFO", message, arguments...)
}

func (l *Input) Warning(message string, arguments ...interface{}) {
	l.Log("WARNING", message, arguments...)
}

func (l *Input) Error(message string, arguments ...interface{}) {
	l.Log("ERROR", message, arguments...)
}

func (l *Input) Critical(message string, arguments ...interface{}) {
	l.Log("CRITICAL", message, arguments...)
}

func (l *Input) Debug(message string, arguments ...interface{}) {
	l.Log("DEBUG", message, arguments...)
}
