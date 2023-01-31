package errors

import "net/http"

func NewMethodNotAllowed(input Input) *Input {
	return New(Input{
		Code:        "MethodNotAllowedError",
		StatusCode:  http.StatusMethodNotAllowed,
		Message:     input.Message,
		SendToSlack: input.SendToSlack,
		Metadata:    input.Metadata,
		Logging:     input.Logging,
	})
}
