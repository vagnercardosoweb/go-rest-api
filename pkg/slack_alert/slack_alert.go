package slack_alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Client struct {
	fields      []Field
	channel     string
	username    string
	memberId    string
	environment string
	color       string
	token       string
}

var (
	pid         = os.Getpid()
	hostname, _ = os.Hostname()
	token       = env.Get("SLACK_TOKEN")
	channel     = env.Get("SLACK_CHANNEL", "alerts")
	username    = env.Get("SLACK_USERNAME", "golang")
	memberId    = env.Get("SLACK_MEMBERS_ID")
)

func NewClient() *Client {
	sa := &Client{
		token:       token,
		channel:     channel,
		username:    username,
		memberId:    memberId,
		environment: config.AppEnv,
		color:       "#D32F2F",
		fields:      make([]Field, 0),
	}
	sa.AddField("Hostname", hostname, true)
	sa.AddField("PID", strconv.Itoa(pid), true)
	return sa
}

func (sa *Client) AddField(title string, value string, short bool) *Client {
	sa.fields = append(sa.fields, Field{Title: title, Value: value, Short: short})
	return sa
}

func (sa *Client) WithError(err *errors.Input) *Client {
	sa.AddField("Error Code / Error Id", fmt.Sprintf("%s / %s", err.Code, err.ErrorId), false)
	sa.AddField("Message", err.Message, false)

	if err.OriginalError != nil {
		sa.AddField(
			"Original Error",
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

func (sa *Client) getColor() string {
	colors := map[string]string{
		"error":   "#D32F2F",
		"warning": "#F57C00",
		"success": "#388E3C",
		"info":    "#0288D1",
	}
	if value, ok := colors[sa.color]; ok {
		return value
	}
	return sa.color
}

func (sa *Client) getMemberIds() string {
	memberIds := strings.Split(sa.memberId, ",")
	if len(memberIds) == 0 {
		return "hey"
	}
	return fmt.Sprintf("<@%s>", strings.Join(memberIds, ">, <@"))
}

func (sa *Client) Send() error {
	if sa.token == "" {
		return nil
	}

	bodyAsBytes, _ := json.Marshal(map[string]any{
		"channel":  sa.channel,
		"username": sa.username,
		"attachments": []map[string]any{{
			"ts": time.Now().UTC().UnixMilli(),
			"text": fmt.Sprintf(
				"%s, an error has occurred",
				sa.getMemberIds(),
			),
			"color":     sa.getColor(),
			"mrkdwn_in": []string{"fields"},
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
			Code:          "SLACK_CREATE_REQUEST",
			Message:       "Slack create request error",
			StatusCode:    http.StatusInternalServerError,
			OriginalError: err.Error(),
			SendToSlack:   errors.Bool(false),
		})
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", sa.token))

	_, err = http.DefaultClient.Do(request)
	if err != nil {
		return errors.New(errors.Input{
			Code:          "SLACK_SEND_REQUEST",
			Message:       "Slack send request error",
			StatusCode:    http.StatusInternalServerError,
			OriginalError: err.Error(),
			SendToSlack:   errors.Bool(false),
		})
	}

	return nil
}
