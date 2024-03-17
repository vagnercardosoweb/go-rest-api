package aws

import (
	"bytes"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
)

type S3Client struct {
	client *s3.S3
	region string
}

func GetS3Client() *S3Client {
	region := env.GetAsString("AWS_S3_REGION", "us-east-1")
	if cached := getServiceFromCache(sesCacheKey, region); cached != nil {
		return cached.(*S3Client)
	}
	client := &S3Client{
		client: s3.New(GetCurrentSession(region)),
		region: region,
	}
	addServiceToCache(sesCacheKey, region, client)
	return client
}

func (c *S3Client) Download(bucket, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}
	output, err := c.client.GetObject(input)
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()
	return io.ReadAll(output.Body)
}

func (c *S3Client) DownloadAndSave(bucket, key, to string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	output, err := c.client.GetObject(input)
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	file, err := os.Create(to)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bodyAsBytes, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	if _, err = io.Copy(file, bytes.NewReader(bodyAsBytes)); err != nil {
		return nil, err
	}

	return bodyAsBytes, nil
}
