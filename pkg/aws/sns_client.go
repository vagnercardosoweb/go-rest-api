package aws

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

type SnsClient struct {
	client *sns.SNS
	region string
}

func GetSnsClient() *SnsClient {
	region := env.GetAsString("AWS_SNS_REGION", "us-east-1")
	if cached := getServiceFromCache(snsCacheKey, region); cached != nil {
		return cached.(*SnsClient)
	}
	client := &SnsClient{
		region: region,
		client: sns.New(GetCurrentSession(region)),
	}
	addServiceToCache(snsCacheKey, region, client)
	return client
}
