package errors

import "net/http"

func NewUnauthorized(input Input) *Input {
	return New(Input{
		Code:        "UnauthorizedError",
		StatusCode:  http.StatusUnauthorized,
		Message:     input.Message,
		SendToSlack: input.SendToSlack,
		Metadata:    input.Metadata,
		Logging:     input.Logging,
	})
}
