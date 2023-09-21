package postgres

import (
	"fmt"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

type Log struct {
	Query        string        `json:"query"`
	Duration     string        `json:"duration"`
	ErrorMessage string        `json:"errorMessage,omitempty"`
	StartedAt    time.Time     `json:"startedAt"`
	FinishedAt   time.Time     `json:"finishedAt"`
	Bind         []interface{} `json:"bind"`
}

func (c *Client) log(log Log) {
	if !c.config.Logging {
		return
	}

	logLevel := logger.LevelInfo
	if log.ErrorMessage != "" {
		logLevel = logger.LevelError
	}

	metadata := map[string]interface{}{
		"tx":         c.tx != nil,
		"query":      log.Query,
		"duration":   log.Duration,
		"startedAt":  log.StartedAt,
		"finishedAt": log.FinishedAt,
		"bind":       log.Bind,
	}

	if log.ErrorMessage != "" {
		metadata["errorMessage"] = log.ErrorMessage
	}

	c.logger.
		WithMetadata(metadata).
		Log(logLevel, fmt.Sprintf("%s_QUERY", c.config.Prefix))

}
