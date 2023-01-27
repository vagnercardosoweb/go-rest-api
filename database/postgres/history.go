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

func (connection *Connection) GetHistory() []History {
	return connection.history
}

func (connection *Connection) addHistory(history History) {
	if connection.hasHistory(history.Query) {
		return
	}
	connection.history = append(connection.history, history)
}

func (connection *Connection) hasHistory(query string) bool {
	result := false

	for _, record := range connection.history {
		if record.Query == query {
			result = true
			break
		}
	}

	return result
}
