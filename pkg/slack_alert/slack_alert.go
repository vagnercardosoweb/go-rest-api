package slack_alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"net/http"
	"strconv"
	"time"
)

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type SlackAlert struct {
	fields      []Field
	channel     string
	username    string
	memberId    string
	environment string
	color       string
	token       string
}

func New() *SlackAlert {
	sa := &SlackAlert{
		token:       env.Get("SLACK_TOKEN"),
		username:    env.Get("SLACK_USERNAME", "go rest api"),
		channel:     env.Get("SLACK_CHANNEL", "golang-alerts"),
		memberId:    env.Get("SLACK_MEMBER_ID"),
		environment: config.AppEnv,
		color:       "#D32F2F",
		fields:      make([]Field, 0),
	}
	sa.AddField("Hostname", config.Hostname, true)
	sa.AddField("PID", strconv.Itoa(config.Pid), true)
	return sa
}

func (sa *SlackAlert) WithToken(token string) *SlackAlert {
	sa.token = token
	return sa
}

func (sa *SlackAlert) WithColor(color string) *SlackAlert {
	sa.color = color
	return sa
}

func (sa *SlackAlert) WithEnvironment(environment string) *SlackAlert {
	sa.environment = environment
	return sa
}

func (sa *SlackAlert) WithUsername(username string) *SlackAlert {
	sa.username = username
	return sa
}

func (sa *SlackAlert) WithChannel(channel string) *SlackAlert {
	sa.channel = channel
	return sa
}

func (sa *SlackAlert) WithMemberId(memberId string) *SlackAlert {
	sa.memberId = memberId
	return sa
}

func (sa *SlackAlert) AddField(title string, value string, short bool) *SlackAlert {
	sa.fields = append(sa.fields, Field{Title: title, Value: value, Short: short})
	return sa
}

func (sa *SlackAlert) WithError(err *errors.Input) *SlackAlert {
	sa.AddField("Error Id", err.ErrorId, false)
	sa.AddField("Error Code", err.Code, false)
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

func (sa *SlackAlert) WithRequestError(method string, path string, err *errors.Input) *SlackAlert {
	sa.AddField("Request - Status", fmt.Sprintf("%s %s - %d", method, path, err.StatusCode), false)
	sa.WithError(err)
	return sa
}

func (sa *SlackAlert) Send() error {
	if sa.token == "" {
		return nil
	}

	bodyAsBytes, _ := json.Marshal(map[string]any{
		"channel":  sa.channel,
		"username": sa.username,
		"attachments": []map[string]any{{
			"ts": time.Now().UTC().UnixMilli(),
			"text": fmt.Sprintf(
				"<@%s>, an error in `%s` using env `%s`",
				sa.memberId,
				sa.username,
				sa.environment,
			),
			"color":     sa.color,
			"mrkdwn_in": []string{"fields"},
			"footer":    sa.username,
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
			Message:       "Slack create Request error",
			StatusCode:    http.StatusInternalServerError,
			OriginalError: err.Error(),
			SendToSlack:   false,
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
			SendToSlack:   false,
		})
	}

	return nil
}
