package errors

import "net/http"

func NewForbidden(input Input) *Input {
	input.Code = "ForbiddenError"
	input.StatusCode = http.StatusForbidden
	return New(input)
}
