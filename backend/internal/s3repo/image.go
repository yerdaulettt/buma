package s3repo

import (
	"bytes"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3Repository struct {
	client *s3.Client
}

func NewS3Repository(c *s3.Client) *s3Repository {
	return &s3Repository{client: c}
}

func (s *s3Repository) PutObject(ctx context.Context, bucket, key string, object []byte) error {
	param := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(object),
	}

	output, err := s.client.PutObject(ctx, param)
	if err != nil {
		return err
	}

	log.Println(output.ResultMetadata)

	return nil
}

func (s *s3Repository) DeleteObject(ctx context.Context, bucket, key string) error {
	param := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	output, err := s.client.DeleteObject(ctx, param)
	if err != nil {
		return err
	}

	log.Println(output.ResultMetadata)

	return nil
}
