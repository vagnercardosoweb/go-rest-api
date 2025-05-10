package mailer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/vagnercardosoweb/go-rest-api/pkg/aws"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
)

type SesClient struct {
	client            *aws.SesClient
	to                []Address
	from              []Address
	replyTo           []Address
	cc                []Address
	bcc               []Address
	configurationName string
	source            string
	template          Template
	subject           string
	html              string
	text              string
	files             []File
}

func NewSesClient() *SesClient {
	return &SesClient{
		client:            aws.GetSesClient(),
		configurationName: env.GetAsString("AWS_SES_CONFIGURATION_NAME"),
		source:            env.Required("AWS_SES_SOURCE"),
	}
}

func (i *SesClient) To(name, address string) Client {
	i.to = append(i.to, Address{Name: name, Address: address})
	return i
}

func (i *SesClient) From(name, address string) Client {
	i.from = append(i.from, Address{Name: name, Address: address})
	return i
}

func (i *SesClient) ReplyTo(name, address string) Client {
	i.replyTo = append(i.replyTo, Address{Name: name, Address: address})
	return i
}

func (i *SesClient) AddCC(name, address string) Client {
	i.cc = append(i.cc, Address{Name: name, Address: address})
	return i
}

func (i *SesClient) AddBCC(name, address string) Client {
	i.bcc = append(i.bcc, Address{Name: name, Address: address})
	return i
}

func (i *SesClient) AddFile(name, path string) Client {
	i.files = append(i.files, File{Name: name, Path: path})
	return i
}

func (i *SesClient) Subject(subject string) Client {
	i.subject = subject
	return i
}

func (i *SesClient) Template(name string, payload any) Client {
	i.template = Template{Name: name, Payload: payload}
	return i
}

func (i *SesClient) Html(value string) Client {
	i.html = value
	return i
}

func (i *SesClient) Text(value string) Client {
	i.text = value
	return i
}

func (i *SesClient) Send() error {
	if len(i.to) == 0 {
		return errors.New(errors.Input{
			Message:    "At least one destination e-mail must be informed.",
			StatusCode: http.StatusBadRequest,
		})
	}

	if i.subject == "" {
		return errors.New(errors.Input{
			Message:    "The subject must be informed to send the email.",
			StatusCode: http.StatusBadRequest,
		})
	}

	if i.template.Name == "" && i.text == "" && i.html == "" {
		return errors.New(errors.Input{
			Message:    "The text or html of the email content needs to be provided",
			StatusCode: http.StatusBadRequest,
		})
	}

	charset := aws.String("UTF-8")
	input := &ses.SendEmailInput{
		Source:           aws.String(i.source),
		ReplyToAddresses: i.parseAddress(i.replyTo),
		Destination: &ses.Destination{
			CcAddresses:  i.parseAddress(i.cc),
			BccAddresses: i.parseAddress(i.bcc),
			ToAddresses:  i.parseAddress(i.to),
		},
	}

	if i.configurationName != "" {
		input.ConfigurationSetName = aws.String(i.configurationName)
	}

	if i.template.Name != "" {
		payloadAsMap := make(map[string]any)

		if value, ok := i.template.Payload.(map[string]any); !ok {
			payloadAsBytes, _ := json.Marshal(i.template.Payload)
			_ = json.Unmarshal(payloadAsBytes, &payloadAsMap)
		} else {
			payloadAsMap = value
		}

		payloadAsMap["year"] = time.Now().Year()
		payloadAsMap["subject"] = i.subject
		payloadAsMap["title"] = i.subject

		payloadAsBytes, _ := json.Marshal(payloadAsMap)

		if _, err := i.client.SendTemplatedEmail(&ses.SendTemplatedEmailInput{
			Source:               input.Source,
			Template:             aws.String(i.template.Name),
			Destination:          input.Destination,
			ConfigurationSetName: input.ConfigurationSetName,
			ReplyToAddresses:     input.ReplyToAddresses,
			TemplateData:         aws.String(string(payloadAsBytes)),
		}); err != nil {
			return err
		}

		return nil
	}

	input.Message = &ses.Message{
		Body: &ses.Body{
			Html: &ses.Content{
				Charset: charset,
				Data:    aws.String(i.html),
			},
			Text: &ses.Content{
				Charset: charset,
				Data:    aws.String(i.text),
			},
		},
		Subject: &ses.Content{
			Charset: charset,
			Data:    aws.String(i.subject),
		},
	}

	if _, err := i.client.SendEmail(input); err != nil {
		return err
	}

	return nil
}

func (*SesClient) parseAddress(addresses []Address) []*string {
	var results []*string

	for _, address := range addresses {
		results = append(
			results,
			aws.String(fmt.Sprintf(
				"%s <%s>",
				address.Name,
				address.Address,
			)),
		)
	}

	return results
}
