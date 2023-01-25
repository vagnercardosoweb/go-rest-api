package errors

import "net/http"

func NewInternalServer(input *AppError) *AppError {
	input.Code = "InternalServerError"
	input.StatusCode = http.StatusInternalServerError
	return New(input)
}
