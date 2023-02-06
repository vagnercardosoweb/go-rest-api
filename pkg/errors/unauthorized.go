package errors

import "net/http"

func NewUnauthorized(input Input) *Input {
	input.Code = "UnauthorizedError"
	input.StatusCode = http.StatusUnauthorized
	return New(input)
}
