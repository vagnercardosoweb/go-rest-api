package aws

import (
	"encoding/json"

	"github.com/vagnercardosoweb/go-rest-api/pkg/env"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/vagnercardosoweb/go-rest-api/pkg/config"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
)

type SqsClient struct {
	client *sqs.SQS
	logger *logger.Logger
	region string
}

func GetSqsClient(logger *logger.Logger) *SqsClient {
	region := env.Get("AWS_SQS_REGION", "us-east-1")
	if cached := getServiceFromCache(sqsCacheKey, region); cached != nil {
		return cached.(*SqsClient)
	}
	client := &SqsClient{
		region: region,
		client: sqs.New(GetCurrentSession(region)),
		logger: logger,
	}
	addServiceToCache(sqsCacheKey, region, client)
	return client
}

func (s *SqsClient) sendMessage(queueUrl string, input any) error {
	if config.IsLocal() {
		return nil
	}
	bodyAsBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}
	s.logger.WithMetadata(map[string]any{"queueUrl": queueUrl, "input": input}).Info("SQS_SEND_MESSAGE_INPUT")
	output, err := s.client.SendMessage(&sqs.SendMessageInput{
		QueueUrl:    String(queueUrl),
		MessageBody: String(string(bodyAsBytes)),
	})
	if err != nil {
		s.logger.AddMetadata("originalError", err.Error()).Info("SQS_SEND_MESSAGE_ERROR")
		return err
	}
	s.logger.AddMetadata("messageId", output.MessageId).Info("SQS_SEND_MESSAGE_COMPLETED")
	return nil
}
