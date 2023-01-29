package errors

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
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

func New(input AppError) *AppError {
	input.checkAndMakeDefaultValues()
	input.AddMetadata("pid", config.Pid)
	input.AddMetadata("hostname", config.Hostname)
	return &input
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) AddMetadata(name string, value interface{}) *AppError {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[name] = value
	return e
}

func (e *AppError) checkAndMakeDefaultValues() {
	if e.Code == "" {
		e.Code = "InternalServerError"
	}
	if e.Message == "" {
		e.Message = "Internal Server Error"
	}
	if e.StatusCode == 0 {
		e.StatusCode = http.StatusInternalServerError
	}
	if e.ErrorId == "" {
		e.ErrorId = fmt.Sprintf("%v", time.Now().UnixMilli())
	}
}
