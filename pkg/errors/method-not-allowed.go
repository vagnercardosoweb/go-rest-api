package errors

import "net/http"

func NewMethodNotAllowed(input Input) *Input {
	input.Code = "MethodNotAllowedError"
	input.StatusCode = http.StatusMethodNotAllowed
	return New(input)
}
