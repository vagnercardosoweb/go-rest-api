package errors

import "net/http"

func NewBadRequest(input Input) *Input {
	return New(Input{
		Code:        "BadRequestError",
		StatusCode:  http.StatusBadRequest,
		Message:     input.Message,
		SendToSlack: input.SendToSlack,
		Metadata:    input.Metadata,
		Logging:     input.Logging,
	})
}
