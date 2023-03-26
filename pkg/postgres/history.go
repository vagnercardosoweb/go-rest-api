package postgres

import "time"

type History struct {
	Query        string        `json:"query"`
	Arguments    []interface{} `json:"arguments"`
	LatencyMs    int64         `json:"latency_ms"`
	ErrorMessage string        `json:"error_message"`
	StartedAt    time.Time     `json:"started_at"`
	FinishedAt   time.Time     `json:"finished_at"`
	CreatedAt    time.Time     `json:"created_at"`
}

func (c *Connection) GetHistory() []History {
	return c.history
}

func (c *Connection) addHistory(history History) {
	if c.hasHistory(history.Query) {
		return
	}
	c.history = append(c.history, history)
}

func (c *Connection) GetLastHistory() History {
	if len(c.history) == 0 {
		return History{}
	}
	return c.history[len(c.history)-1]

}

func (c *Connection) hasHistory(query string) bool {
	result := false

	for _, record := range c.history {
		if record.Query == query {
			result = true
			break
		}
	}

	return result
}
