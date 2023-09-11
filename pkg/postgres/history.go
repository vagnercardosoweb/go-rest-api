package postgres

import (
	"fmt"
	"reflect"
	"time"
)

type History struct {
	Query        string        `json:"query"`
	Duration     string        `json:"duration"`
	ErrorMessage string        `json:"errorMessage"`
	StartedAt    time.Time     `json:"startedAt"`
	FinishedAt   time.Time     `json:"finishedAt"`
	Bind         []interface{} `json:"bind"`
}

func (c *Client) GetHistory() []History {
	return c.history
}

func (c *Client) addHistory(history History) {
	if c.config.Logging {
		c.logger.
			WithMetadata(map[string]interface{}{
				"query":         history.Query,
				"inTransaction": c.tx != nil,
				"errorMessage":  history.ErrorMessage,
				"duration":      history.Duration,
				"bind":          history.Bind,
			}).
			Info(fmt.Sprintf("%s_QUERY", c.config.Prefix))
	}
	if c.hasHistory(history) {
		return
	}
	c.history = append(c.history, history)
}

func (c *Client) GetLastHistory() History {
	if len(c.history) == 0 {
		return History{}
	}
	return c.history[len(c.history)-1]

}

func (c *Client) hasHistory(history History) bool {
	result := false

	for _, record := range c.history {
		if record.Query == history.Query && reflect.DeepEqual(record.Bind, history.Bind) {
			result = true
			break
		}
	}

	return result
}
