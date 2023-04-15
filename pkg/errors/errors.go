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
		SendToSlack bool     `json:"sendToSlack"`
		Arguments   []any    `json:"arguments"`
	}
)

func New(input Input) *Input {
	input.validate()

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
	if input.StatusCode == 0 {
		input.StatusCode = http.StatusInternalServerError
	}
	if input.Message == "" {
		input.Message = http.StatusText(input.StatusCode)
	}
	if input.ErrorId == "" {
		input.ErrorId = uuid.New().String()
	}
	if input.Code == "" {
		input.Code = "DEFAULT"
	}
	input.Message = fmt.Sprintf(input.Message, input.Arguments...)
}
