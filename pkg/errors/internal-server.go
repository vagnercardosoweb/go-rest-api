package errors

import "net/http"

func NewInternalServer(input Input) *Input {
	return New(Input{
		Code:        "InternalServerError",
		StatusCode:  http.StatusInternalServerError,
		Message:     input.Message,
		SendToSlack: input.SendToSlack,
		Metadata:    input.Metadata,
		Logging:     input.Logging,
	})
}
