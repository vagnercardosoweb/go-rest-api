package aws

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vagnercardosoweb/go-rest-api/pkg/env"
	"github.com/vagnercardosoweb/go-rest-api/pkg/logger"
	"github.com/vagnercardosoweb/go-rest-api/pkg/utils"
)

type S3Client struct {
	ctx    context.Context
	client *s3.Client
	region string
}

func GetS3Client(ctx context.Context, logger *logger.Logger) *S3Client {
	region := env.GetAsString("AWS_S3_REGION", "us-east-1")

	if cached := getServiceFromCache(s3CacheKey, region); cached != nil {
		return cached.(*S3Client)
	}

	cfg, _ := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	client := &S3Client{ctx: ctx, client: s3.NewFromConfig(cfg), region: region}

	addServiceToCache(s3CacheKey, region, client)

	return client
}

func (c *S3Client) DownloadAsBytes(bucket, key string) ([]byte, error) {
	input := &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	}

	output, err := c.client.GetObject(c.ctx, input)
	if err != nil {
		return nil, err
	}

	defer output.Body.Close()

	return io.ReadAll(output.Body)
}

func (c *S3Client) DownloadAsFile(bucket, key, to string) ([]byte, error) {
	if err := utils.ValidateDownloadFilePath(to); err != nil {
		return nil, err
	}

	dir := filepath.Dir(to)
	filename := filepath.Base(to)
	sanitizedFilename := utils.SanitizeFileName(filename)
	safePath := filepath.Join(dir, sanitizedFilename)

	input := &s3.GetObjectInput{Bucket: &bucket, Key: &key}
	output, err := c.client.GetObject(c.ctx, input)
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	// #nosec G301 - Directory permissions are intentionally 0750 for security
	if err := os.MkdirAll(dir, 0750); err != nil {
		return nil, err
	}

	// #nosec G304 - Path is validated by ValidateDownloadFilePath above
	file, err := os.Create(safePath)
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

func (c *S3Client) PutSignedURL(bucket, key string, expiresIn time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.client)

	input := &s3.PutObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)}
	req, err := presignClient.PresignPutObject(c.ctx, input, s3.WithPresignExpires(expiresIn))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}

func (c *S3Client) GetSignedURL(bucket, key string, expiresIn time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(c.client)

	input := &s3.GetObjectInput{Bucket: aws.String(bucket), Key: aws.String(key)}
	req, err := presignClient.PresignGetObject(c.ctx, input, s3.WithPresignExpires(expiresIn))

	if err != nil {
		return "", err
	}

	return req.URL, nil
}
