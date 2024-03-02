package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

var hostname, _ = os.Hostname()

func NewAlert() *Client {
	return &Client{
		token:       env.GetAsString("SLACK_TOKEN"),
		channel:     env.GetAsString("SLACK_CHANNEL", "alerts"),
		username:    env.GetAsString("SLACK_USERNAME", "golang"),
		isEnabled:   env.GetAsBool("SLACK_ENABLED", "true"),
		memberIds:   strings.Split(env.GetAsString("SLACK_MEMBERS_ID"), ","),
		environment: env.GetAppEnv(),
		color:       ColorError,
		fields:      make([]Field, 0),
		mu:          &sync.Mutex{},
	}
}

func (sa *Client) AddField(title string, value string, short bool) *Client {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	sa.fields = append(sa.fields, Field{Title: title, Value: value, Short: short})
	return sa
}

func (sa *Client) AddError(title string, err any) *Client {
	sa.AddField(title, fmt.Sprintf("```%s```", err), false)
	return sa
}

func (sa *Client) WithError(err *errors.Input) *Client {
	sa.AddField("ErrorCode / ErrorId", fmt.Sprintf("%s / %s", err.Code, err.ErrorId), false)
	sa.AddField("Message", err.Message, false)

	if err.OriginalError != nil {
		sa.AddField(
			"Error",
			fmt.Sprintf("```%s```", err.OriginalError),
			false,
		)
	}

	return sa
}

func (sa *Client) WithRequestError(method string, path string, err *errors.Input) *Client {
	sa.AddField("[Status] Request", fmt.Sprintf("[%d] %s %s", err.StatusCode, method, path), false)
	sa.WithError(err)
	return sa
}

func (sa *Client) WithColor(color ColorName) *Client {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	sa.color = color
	return sa
}

func (sa *Client) getColor() string {
	errorColor := "#D32F2F"
	colors := map[ColorName]string{
		"error":   errorColor,
		"warning": "#F57C00",
		"success": "#388E3C",
		"info":    "#0288D1",
	}
	if value, ok := colors[sa.color]; ok {
		return value
	}
	return errorColor
}

func (sa *Client) WithMemberId(memberId string) *Client {
	sa.mu.Lock()
	defer sa.mu.Unlock()
	sa.memberIds = append(sa.memberIds, memberId)
	return sa
}

func (sa *Client) getMemberIds() string {
	if len(sa.memberIds) == 0 {
		return "hey"
	}
	return fmt.Sprintf("<@%s>", strings.Join(sa.memberIds, ">, <@"))
}

func (sa *Client) Send() error {
	if !sa.isEnabled || sa.token == "" {
		return nil
	}

	bodyAsBytes, _ := json.Marshal(map[string]any{
		"channel":  sa.channel,
		"username": sa.username,
		"attachments": []map[string]any{{
			"ts": time.Now().UTC().UnixMilli(),
			"text": fmt.Sprintf(
				"%s, new alert in `%s`.",
				sa.getMemberIds(),
				hostname,
			),
			"color":     sa.getColor(),
			"mrkdwn_in": []string{"text", "fields"},
			"footer":    fmt.Sprintf("[%s] %s", sa.environment, sa.username),
			"fields":    sa.fields,
		}},
	})

	request, err := http.NewRequest(
		http.MethodPost,
		"https://slack.com/api/chat.postMessage",
		bytes.NewBuffer(bodyAsBytes),
	)

	if err != nil {
		return errors.New(errors.Input{
			Code:        "SLACK_NEW_REQUEST",
			SendToSlack: errors.Bool(false),
			Message:     err.Error(),
		})
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sa.token))

	_, err = http.DefaultClient.Do(request)
	if err != nil {
		return errors.New(errors.Input{
			Code:        "SLACK_SEND_REQUEST",
			SendToSlack: errors.Bool(false),
			Message:     err.Error(),
		})
	}

	return nil
}
