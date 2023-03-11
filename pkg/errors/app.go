package errors

import (
	"fmt"
	"net/http"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
)

type (
	Metadata map[string]any
	Input    struct {
		Code        string   `json:"code"`
		ErrorId     string   `json:"errorId"`
		Message     string   `json:"message"`
		StatusCode  int      `json:"statusCode"`
		Metadata    Metadata `json:"metadata"`
		Logging     bool     `json:"-"`
		Arguments   []any    `json:"-"`
		SendToSlack bool     `json:"-"`
	}
)

func New(input Input) *Input {
	input.makeDefaultValues()
	input.AddMetadata("pid", config.Pid)
	input.AddMetadata("hostname", config.Hostname)
	return &input
}

func (e *Input) Error() string {
	return e.Message
}

func (e *Input) AddMetadata(name string, value any) *Input {
	if e.Metadata == nil {
		e.Metadata = make(Metadata)
	}
	e.Metadata[name] = value
	return e
}

func (e *Input) makeDefaultValues() {
	if e.Code == "" {
		e.Code = "InternalServerError"
	}
	if e.Message == "" {
		e.Message = "InternalServerError"
	}
	if e.StatusCode == 0 {
		e.StatusCode = http.StatusInternalServerError
	}
	if e.ErrorId == "" {
		e.ErrorId = fmt.Sprintf("%v", time.Now().UnixMilli())
	}
	e.Message = fmt.Sprintf(e.Message, e.Arguments...)
}
