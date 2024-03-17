package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func New(input Input) *Input {
	input.build()
	input.makeStack()
	return &input
}

func NewFromString(message string) *Input {
	return New(Input{Message: message})
}

// FromSql converts a sql error to an AppError.
// First argument is the error, the rest are arguments to be used in the message
func FromSql(err error, args ...any) *Input {
	appError := New(Input{OriginalError: err})

	if Is(err, sql.ErrNoRows) {
		appError.Message = "Resource not found"
		appError.StatusCode = http.StatusNotFound
		appError.SendToSlack = Bool(false)
	}

	if len(args) > 0 {
		appError.Message = args[0].(string)
		appError.Arguments = args[1:]
		appError.makeMessage()
	}

	return appError
}

func FromBindJson(err error, translator ut.Translator) *Input {
	appError := New(Input{
		Message:    err.Error(),
		StatusCode: http.StatusUnprocessableEntity,
		Code:       "BIND_JSON_ERROR",
	})

	if Is(err, io.EOF) {
		appError.Message = "Error retrieving the request body"
		appError.OriginalError = err.Error()
	}

	if translator == nil {
		return appError
	}

	var errs validator.ValidationErrors
	if errors.As(err, &errs) {
		validations := make([]map[string]any, 0)

		appError.Message = "Some fields are invalid"
		appError.Code = "VALIDATION_ERROR"

		for _, e := range errs {
			validations = append(validations, map[string]any{
				"tag":       e.Tag(),
				"field":     e.Field(),
				"message":   e.Translate(translator),
				"namespace": e.Namespace(),
				"value":     e.Value(),
				"param":     e.Param(),
			})
		}

		appError.Message = validations[0]["message"].(string)
		appError.Metadata["validations"] = validations
	}

	return appError
}

func (e *Input) Error() string {
	return e.Message
}

func (e *Input) makeStack() {
	if e.SkipStack == true {
		return
	}

	if len(e.Stack) == 0 {
		e.Stack = GetCallerStack(3)
	}
}

func (e *Input) checkInputValues() {
	if e.Name == "" {
		e.Name = "AppError"
	}

	if e.Code == "" {
		e.Code = "DEFAULT"
	}

	if e.StatusCode == 0 {
		e.StatusCode = http.StatusInternalServerError
	}

	if e.Logging == nil {
		e.Logging = Bool(true)
	}
}

func (e *Input) build() {
	e.makeMetadata()
	e.checkInputValues()
	e.checkSendToSlack()
	e.checkOriginalError()
	e.makeMessage()
}

func (e *Input) checkSendToSlack() {
	if e.SendToSlack != nil {
		return
	}

	switch e.StatusCode {
	case
		http.StatusNotFound,
		http.StatusForbidden,
		http.StatusUnprocessableEntity,
		http.StatusBadRequest,
		http.StatusUnauthorized:
		e.SendToSlack = Bool(false)
	default:
		e.SendToSlack = Bool(true)
	}
}

func (e *Input) makeMessage() {
	e.Message = fmt.Sprintf(e.Message, e.Arguments...)
	if e.Message == "" {
		e.Message = http.StatusText(e.StatusCode)
	}
}

func (e *Input) checkOriginalError() {
	if _, ok := e.OriginalError.(*Input); ok {
		return
	}

	if err, ok := e.OriginalError.(error); ok {
		e.OriginalError = err.Error()
	}
}

func (e *Input) makeMetadata() {
	if e.Metadata == nil {
		e.Metadata = make(Metadata)
	}

	for name, value := range e.Metadata {
		if _, ok := value.(*Input); ok {
			continue
		}

		if err, ok := value.(error); ok {
			e.Metadata[name] = err.Error()
		}
	}
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Bool(value bool) *bool {
	return &value
}
