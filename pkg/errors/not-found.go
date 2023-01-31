package errors

import "net/http"

func NewNotFound(input Input) *Input {
	return New(Input{
		Code:        "NotFoundError",
		StatusCode:  http.StatusNotFound,
		Message:     input.Message,
		SendToSlack: input.SendToSlack,
		Metadata:    input.Metadata,
		Logging:     input.Logging,
	})
}
