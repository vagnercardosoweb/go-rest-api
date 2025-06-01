package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

type SnsClient struct {
	ctx    context.Context
	client *sns.Client
	region string
}

func GetSnsClient(ctx context.Context) *SnsClient {
	region := env.GetAsString("AWS_REGION", "us-east-1")

	if cached := getServiceFromCache(snsCacheKey, region); cached != nil {
		return cached.(*SnsClient)
	}

	cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	client := &SnsClient{ctx: ctx, client: sns.NewFromConfig(cfg), region: region}

	addServiceToCache(snsCacheKey, region, client)

	return client
}
