package errors

import "net/http"

func NewNotFound(message string, metadata map[string]interface{}) *AppError {
	return New(&AppError{
		Code:       "NotFoundError",
		StatusCode: http.StatusNotFound,
		Message:    message,
		Metadata:   metadata,
	})
}
