package errors

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
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
	input.validate()
	input.defineCode()

	if input.Metadata == nil {
		input.Metadata = Metadata{
			"pid":      config.Pid,
			"hostname": config.Hostname,
		}
	}

	return &input
}

func (e *Input) Error() string {
	return e.Message
}

func (input *Input) AddMetadata(name string, value any) *Input {
	input.Metadata[name] = value
	return input
}

func (input *Input) validate() {
	if input.Message == "" {
		input.Message = "InternalServerError"
	}
	if input.StatusCode == 0 {
		input.StatusCode = http.StatusInternalServerError
	}
	if input.ErrorId == "" {
		input.ErrorId = uuid.New().String()
	}
	input.Message = fmt.Sprintf(input.Message, input.Arguments...)
}

func (input *Input) defineCode() {
	if input.Code == "" {
		input.Code = http.StatusText(input.StatusCode)
	}
}
