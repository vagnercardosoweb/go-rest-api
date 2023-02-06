package errors

import "net/http"

func NewInternalServer(input Input) *Input {
	input.Code = "InternalServerError"
	input.StatusCode = http.StatusInternalServerError
	return New(input)
}
