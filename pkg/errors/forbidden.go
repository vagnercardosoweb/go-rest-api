package errors

import "net/http"

func NewForbidden(input Input) *Input {
	return New(Input{
		Code:        "ForbiddenError",
		StatusCode:  http.StatusForbidden,
		Message:     input.Message,
		SendToSlack: input.SendToSlack,
		Metadata:    input.Metadata,
		Logging:     input.Logging,
	})
}
