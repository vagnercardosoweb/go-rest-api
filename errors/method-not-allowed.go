package errors

import "net/http"

func NewMethodNotAllowed(message string, metadata map[string]interface{}) *AppError {
	return New(&AppError{
		Code:       "MethodNotAllowedError",
		StatusCode: http.StatusMethodNotAllowed,
		Message:    message,
		Metadata:   metadata,
	})
}
