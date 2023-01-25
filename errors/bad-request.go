package errors

import "net/http"

func NewBadRequest(message string) *AppError {
	return New(&AppError{
		Code:       "BadRequestError",
		StatusCode: http.StatusBadRequest,
		Message:    message,
	})
}
