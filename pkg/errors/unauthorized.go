package errors

import "net/http"

func NewUnauthorized(message string) *AppError {
	return New(AppError{
		Code:       "UnauthorizedError",
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	})
}
