package postgres

import (
	"regexp"
	"strings"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

type Log struct {
	Query        string    `json:"query"`
	Duration     string    `json:"duration"`
	ErrorMessage string    `json:"errorMessage,omitempty"`
	FinishedAt   time.Time `json:"finishedAt"`
	StartedAt    time.Time `json:"startedAt"`
	Bind         []any     `json:"bind"`
}

func (c *Client) log(log *Log) {
	log.Duration = log.FinishedAt.Sub(log.StartedAt).String()
	c.lastLog = log

	if !c.config.Logging {
		return
	}

	logLevel := logger.LevelInfo
	metadata := map[string]any{
		"tx":         c.tx != nil,
		"query":      log.getQuery(),
		"startedAt":  log.StartedAt,
		"finishedAt": log.FinishedAt,
		"duration":   log.Duration,
		"bind":       log.Bind,
	}

	if log.ErrorMessage != "" {
		metadata["errorMessage"] = log.ErrorMessage
		logLevel = logger.LevelError
	}

	c.logger.
		WithFields(metadata).
		Log(logLevel, "DB_QUERY")
}

func (l *Log) getQuery() string {
	q := strings.TrimSpace(l.Query)

	q = regexp.MustCompile(`\s+|\n|\t`).ReplaceAllString(q, " ")
	q = regexp.MustCompile(`\(\s`).ReplaceAllString(q, "(")
	q = regexp.MustCompile(`\s\)`).ReplaceAllString(q, ")")

	return q
}
