package errors

import "net/http"

func NewForbidden(message string) *AppError {
	return New(&AppError{
		Code:       "ForbiddenError",
		StatusCode: http.StatusForbidden,
		Message:    message,
	})
}
