package errors

import (
	"fmt"
	"net/http"
	"rest-api/config"
	"time"
)

type AppError struct {
	Code        string                 `json:"code"`
	Message     string                 `json:"message"`
	ErrorId     string                 `json:"errorId"`
	StatusCode  int                    `json:"statusCode"`
	Metadata    map[string]interface{} `json:"metadata"`
	SendToSlack bool                   `json:"-"`
	Logging     bool                   `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func New(input *AppError) *AppError {
	if input.StatusCode == 0 {
		input.StatusCode = http.StatusInternalServerError
	}
	if input.Code == "" {
		input.Code = "InternalServerError"
	}
	if input.Message == "" {
		input.Message = "Internal Server Error"
	}
	if len(input.Metadata) == 0 {
		input.Metadata = map[string]interface{}{}
	}
	if input.ErrorId == "" {
		input.ErrorId = fmt.Sprintf("%d", time.Now().UnixMilli())
	}
	input.AddMetadata("pid", config.Pid)
	input.AddMetadata("hostname", config.Hostname)
	return input
}

func (e *AppError) AddMetadata(name string, value interface{}) *AppError {
	e.Metadata[name] = value
	return e
}
