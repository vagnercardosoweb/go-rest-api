package aws

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	awsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

type SqsClient struct {
	ctx    context.Context
	client *sqs.Client
	logger *logger.Logger
	region string
}

func GetSqsClient(ctx context.Context, logger *logger.Logger) *SqsClient {
	region := env.GetAsString("AWS_REGION", "us-east-1")

	if cached := getServiceFromCache(sqsCacheKey, region); cached != nil {
		return cached.(*SqsClient)
	}

	cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	client := &SqsClient{ctx: ctx, region: region, client: sqs.NewFromConfig(cfg), logger: logger}

	addServiceToCache(sqsCacheKey, region, client)

	return client
}

func (s *SqsClient) SendMessage(queueUrl *string, input any) error {
	if env.IsLocal() {
		s.logger.
			AddMetadata("url", queueUrl).
			AddMetadata("input", input).
			Info("SQS_SEND_MESSAGE_LOCAL")

		return nil
	}

	inputAsBytes, _ := json.Marshal(input)
	s.logger.
		AddMetadata("url", queueUrl).
		AddMetadata("input", input).
		Info("SQS_SEND_MESSAGE_INPUT")

	sendMessageInput := &sqs.SendMessageInput{
		QueueUrl:    queueUrl,
		MessageBody: String(string(inputAsBytes)),
		MessageAttributes: map[string]awsTypes.MessageAttributeValue{
			"SendFrom": {
				StringValue: String("go-rest-api"),
				DataType:    String("String"),
			},
		},
	}

	if strings.HasSuffix(*queueUrl, ".fifo") {
		sendMessageInput.MessageGroupId = String("default")
	}

	output, err := s.client.SendMessage(s.ctx, sendMessageInput)
	if err != nil {
		s.logger.
			AddMetadata("error", err.Error()).
			Info("SQS_SEND_MESSAGE_ERROR")

		return err
	}

	s.logger.
		AddMetadata("messageId", output.MessageId).
		Info("SQS_SEND_MESSAGE_COMPLETED")

	return nil
}
