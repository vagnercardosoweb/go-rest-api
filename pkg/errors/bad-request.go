package errors

import "net/http"

func NewBadRequest(input Input) *Input {
	input.Code = "BadRequestError"
	input.StatusCode = http.StatusBadRequest
	return New(input)
}
