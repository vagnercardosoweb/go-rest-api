package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type (
	Metadata map[string]any
	Input    struct {
		Name          string   `json:"name"`
		Code          string   `json:"code"`
		ErrorId       string   `json:"errorId"`
		Message       string   `json:"message"`
		StatusCode    int      `json:"statusCode"`
		SendToSlack   *bool    `json:"sendToSlack"`
		Logging       *bool    `json:"logging"`
		Metadata      Metadata `json:"metadata,omitempty"`
		Arguments     []any    `json:"-"`
		OriginalError any      `json:"originalError,omitempty"`
		SkipStack     bool     `json:"-"`
		Stack         []string `json:"stack"`
	}
)

var translator ut.Translator

func init() {
	if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
		lang := en.New()
		translator, _ = ut.New(lang, lang).GetTranslator("en")
		_ = enTranslation.RegisterDefaultTranslations(val, translator)
	}
}

func New(input Input) *Input {
	input.build()
	input.makeStack()
	return &input
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Bool(value bool) *bool {
	return &value
}

// FromSql converts a sql error to an AppError.
// First argument is the error, the rest are arguments to be used in the message
func FromSql(err error, args ...any) *Input {
	appError := New(Input{OriginalError: err})

	if Is(err, sql.ErrNoRows) {
		appError.Message = "Resource not found"
		appError.StatusCode = http.StatusNotFound

		falsy := Bool(false)
		appError.SendToSlack = falsy
		appError.Logging = falsy
	}

	if len(args) > 0 {
		appError.Message = args[0].(string)
		appError.Arguments = args[1:]
		appError.makeMessage()
	}

	return appError
}

func FromBindJson(err error) *Input {
	falsy := Bool(false)
	appError := New(Input{
		Code:        "BIND_JSON_ERROR",
		Message:     err.Error(),
		StatusCode:  http.StatusUnprocessableEntity,
		SendToSlack: falsy,
		Logging:     falsy,
	})

	if Is(err, io.EOF) {
		appError.Message = "Error retrieving the request body, please check that the data is correct."
		appError.OriginalError = err.Error()
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
		e.Stack = GetCallerStack(2)
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

	truthy := Bool(true)
	if e.Logging == nil {
		e.Logging = truthy
	}

	if e.SendToSlack == nil {
		e.SendToSlack = truthy
	}
}

func (e *Input) build() {
	e.makeMetadata()
	e.checkInputValues()
	e.checkOriginalError()
	e.makeMessage()
}

func (e *Input) makeMessage() {
	e.Message = fmt.Sprintf(
		e.Message,
		e.Arguments...,
	)

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
