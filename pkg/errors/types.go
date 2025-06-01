package errors

type (
	Metadata map[string]any
	Input    struct {
		Code          string   `json:"code,omitempty"`
		Name          string   `json:"name,omitempty"`
		Message       string   `json:"message,omitempty"`
		StatusCode    int      `json:"statusCode,omitempty"`
		SendToSlack   *bool    `json:"-"`
		RequestId     string   `json:"-"`
		Logging       *bool    `json:"-"`
		Metadata      Metadata `json:"metadata,omitempty"`
		Arguments     []any    `json:"-"`
		OriginalError any      `json:"originalError,omitempty"`
		SkipStack     bool     `json:"-"`
		Stack         []string `json:"stack,omitempty"`
	}
)
