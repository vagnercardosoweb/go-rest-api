package mailer

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/errors"
	"net/http"
	"time"
)

type SES struct {
	client            *ses.SES
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

func NewSes() Mailer {
	source := env.Required("AWS_SES_SOURCE")
	configurationName := env.Get("AWS_SES_CONFIGURATION_NAME")
	region := env.Get("AWS_SES_REGION", "us-east-1")
	client := ses.New(
		session.Must(session.NewSession()),
		aws.NewConfig().WithRegion(region),
	)
	return &SES{
		client:            client,
		configurationName: configurationName,
		source:            source,
	}
}

func (i *SES) To(name, address string) Mailer {
	i.to = append(i.to, Address{Name: name, Address: address})
	return i
}

func (i *SES) From(name, address string) Mailer {
	i.from = append(i.from, Address{Name: name, Address: address})
	return i
}

func (i *SES) ReplyTo(name, address string) Mailer {
	i.replyTo = append(i.replyTo, Address{Name: name, Address: address})
	return i
}

func (i *SES) AddCC(name, address string) Mailer {
	i.cc = append(i.cc, Address{Name: name, Address: address})
	return i
}

func (i *SES) AddBCC(name, address string) Mailer {
	i.bcc = append(i.bcc, Address{Name: name, Address: address})
	return i
}

func (i *SES) AddFile(name, path string) Mailer {
	i.files = append(i.files, File{Name: name, Path: path})
	return i
}

func (i *SES) Subject(subject string) Mailer {
	i.subject = subject
	return i
}

func (i *SES) Template(name string, payload TemplatePayload) Mailer {
	i.template = Template{Name: name, Payload: payload}
	return i
}

func (i *SES) Html(value string) Mailer {
	i.html = value
	return i
}

func (i *SES) Text(value string) Mailer {
	i.text = value
	return i
}

func (i *SES) Send() error {
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
		if &i.template.Payload == nil {
			i.template.Payload = make(TemplatePayload)
		}
		i.template.Payload["year"] = time.Now().Year()
		i.template.Payload["subject"] = i.subject
		if _, ok := i.template.Payload["title"]; !ok {
			i.template.Payload["title"] = i.subject
		}
		payloadAsBytes, _ := json.Marshal(i.template.Payload)
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

func (*SES) parseAddress(addresses []Address) []*string {
	var results []*string
	for _, address := range addresses {
		results = append(
			results,
			aws.String(fmt.Sprintf("%s <%s>", address.Name, address.Address)),
		)
	}
	return results
}
