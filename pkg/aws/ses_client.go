package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

type SesClient struct {
	ctx    context.Context
	client *ses.Client
	region string
}

func GetSesClient(ctx context.Context) *SesClient {
	region := env.GetAsString("AWS_REGION", "us-east-1")

	if cached := getServiceFromCache(sesCacheKey, region); cached != nil {
		return cached.(*SesClient)
	}

	cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	client := &SesClient{ctx: ctx, client: ses.NewFromConfig(cfg), region: region}

	addServiceToCache(sesCacheKey, region, client)

	return client
}

func (s *SesClient) SendEmailWithTemplate(input *ses.SendTemplatedEmailInput) (*ses.SendTemplatedEmailOutput, error) {
	return s.client.SendTemplatedEmail(s.ctx, input)
}

func (s *SesClient) SendEmail(input *ses.SendEmailInput) (*ses.SendEmailOutput, error) {
	return s.client.SendEmail(s.ctx, input)
}
