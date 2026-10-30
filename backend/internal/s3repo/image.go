package s3repo

import (
	"bytes"
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type s3Repository struct {
	client    *s3.Client
	presigner *s3.PresignClient
}

func NewS3Repository(c *s3.Client) *s3Repository {
	return &s3Repository{client: c, presigner: s3.NewPresignClient(c)}
}

func (s *s3Repository) PutObject(ctx context.Context, bucket, key string, object []byte) error {
	param := &s3.PutObjectInput{
		Bucket:               aws.String(bucket),
		Key:                  aws.String(key),
		Body:                 bytes.NewReader(object),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	}

	output, err := s.client.PutObject(ctx, param)
	if err != nil {
		return err
	}

	log.Println(output.ResultMetadata)

	return nil
}

func (s *s3Repository) GetObjectUrl(ctx context.Context, bucket, key string, limit time.Duration) (string, error) {
	param := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := s.presigner.PresignGetObject(ctx, param, func(po *s3.PresignOptions) {
		po.Expires = limit
	})
	if err != nil {
		return "", err
	}

	s3url := strings.TrimLeft(result.URL, os.Getenv("AWS_ENDPOINT_URL"))

	return s3url, nil
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
