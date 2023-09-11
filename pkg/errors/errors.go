package errors

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslation "github.com/go-playground/validator/v10/translations/en"
	"github.com/google/uuid"
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
		SendToSlack   *bool    `json:"sendToSlack"`
		Arguments     []any    `json:"arguments"`
		Logging       *bool    `json:"logging"`
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

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func Bool(value bool) *bool {
	return &value
}

func FromSql(err error, errorMessage ...string) *Input {
	appError := New(Input{})
	appError.OriginalError = err.Error()
	appError.StatusCode = http.StatusInternalServerError

	if Is(err, sql.ErrNoRows) {
		appError.Message = "Resource not found"
		appError.StatusCode = http.StatusNotFound
		falsy := Bool(false)
		appError.SendToSlack = falsy
		appError.Logging = falsy
	}

	if len(errorMessage) > 0 {
		appError.Message = errorMessage[0]
	}

	return appError
}

func FromBindJson(err error) *Input {
	falsy := Bool(false)
	appError := New(Input{
		Code:        "BIND_JSON_ERROR",
		Message:     err.Error(),
		StatusCode:  http.StatusBadRequest,
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

		if val, ok := binding.Validator.Engine().(*validator.Validate); ok {
			lang := en.New()
			trans, _ := ut.New(lang, lang).GetTranslator("en")
			_ = enTranslation.RegisterDefaultTranslations(val, trans)

			for _, e := range errs {
				validations = append(validations, map[string]any{
					"tag":       e.Tag(),
					"field":     e.Field(),
					"message":   e.Translate(trans),
					"namespace": e.Namespace(),
					"value":     e.Value(),
					"param":     e.Param(),
				})
			}

			appError.Message = validations[0]["message"].(string)
			_ = appError.AddMetadata("validations", validations)
		}
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

	truthy := Bool(true)
	if input.Logging == nil {
		input.Logging = truthy
	}

	if input.SendToSlack == nil {
		input.SendToSlack = truthy
	}

	input.Message = fmt.Sprintf(
		input.Message,
		input.Arguments...,
	)
}
