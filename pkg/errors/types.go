package errors

type (
	Metadata map[string]any
	Input    struct {
		Name          string   `json:"name"`
		Code          string   `json:"code"`
		ErrorId       string   `json:"errorId"`
		Message       string   `json:"message"`
		StatusCode    int      `json:"statusCode"`
		SendToSlack   *bool    `json:"sendToSlack"`
		Logging       *bool    `json:"logging"`
		Metadata      Metadata `json:"metadata,omitempty"`
		Arguments     []any    `json:"-"`
		OriginalError any      `json:"originalError,omitempty"`
		SkipStack     bool     `json:"-"`
		Stack         []string `json:"stack"`
	}
)
