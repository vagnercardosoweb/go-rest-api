package errors

import "net/http"

func NewNotFound(input Input) *Input {
	input.Code = "NotFoundError"
	input.StatusCode = http.StatusNotFound
	return New(input)
}
