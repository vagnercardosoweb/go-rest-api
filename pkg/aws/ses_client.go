package aws

import (
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

type SesClient struct {
	client *ses.SES
	region string
}

func GetSesClient() *SesClient {
	region := env.GetAsString("AWS_SES_REGION", "us-east-1")
	if cached := getServiceFromCache(sesCacheKey, region); cached != nil {
		return cached.(*SesClient)
	}
	client := &SesClient{
		client: ses.New(GetCurrentSession(region)),
		region: region,
	}
	addServiceToCache(sesCacheKey, region, client)
	return client
}

func (s *SesClient) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	return s.client.SendEmail(input)
}

func (s *SesClient) SendTemplatedEmail(input *ses.SendTemplatedEmailInput) (*ses.SendTemplatedEmailOutput, error) {
	return s.client.SendTemplatedEmail(input)
}
