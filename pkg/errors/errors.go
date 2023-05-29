package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
)

type (
	Metadata map[string]any
	Input    struct {
		Name          string   `json:"name"`
		Code          string   `json:"code"`
		ErrorId       string   `json:"errorId"`
		Message       string   `json:"message"`
		StatusCode    int      `json:"statusCode"`
		OriginalError any      `json:"originalError"`
		Stack         []string `json:"stack"`
		Metadata      Metadata `json:"metadata"`
		SendToSlack   bool     `json:"sendToSlack"`
		Arguments     []any    `json:"arguments"`
	}
)

func New(input Input) *Input {
	input.makeDefaultValues()
	if len(input.Stack) == 0 {
		input.Stack = GetCallerStack(2)
	}
	return &input
}

func As(err error, target any) bool {
	return errors.As(err, &target)
}

func WithSqlError(err error, errorMessage ...string) *Input {
	appError := New(Input{})
	appError.OriginalError = err.Error()
	appError.StatusCode = http.StatusInternalServerError

	if err == sql.ErrNoRows {
		appError.Message = "No rows in result set"
		appError.StatusCode = http.StatusNotFound
	} else {
		appError.SendToSlack = true
	}

	if len(errorMessage) > 0 {
		appError.Message = errorMessage[0]
	}

	return appError
}

func (input *Input) Error() string {
	return input.Message
}

func (input *Input) AddMetadata(name string, value any) *Input {
	input.Metadata[name] = value
	return input
}

func (input *Input) makeDefaultValues() {
	if originalError, ok := input.OriginalError.(*Input); ok {
		*input = *originalError
	} else if err, ok := input.OriginalError.(error); ok {
		input.OriginalError = err.Error()
	}

	if input.Name == "" {
		input.Name = "AppError"
	}

	if input.Code == "" {
		input.Code = "DEFAULT"
	}

	if input.StatusCode == 0 {
		input.StatusCode = http.StatusInternalServerError
	}

	if input.Message == "" {
		input.Message = http.StatusText(input.StatusCode)
	}

	if input.ErrorId == "" {
		input.ErrorId = uuid.New().String()
	}

	if input.Metadata == nil {
		input.Metadata = make(Metadata)
	}

	input.Message = fmt.Sprintf(
		input.Message,
		input.Arguments...,
	)
}
