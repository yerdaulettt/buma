package image

import (
	"context"
	"time"
)

type imageRepository interface {
	AddMetadata(ctx context.Context, userId int, imageType, key string) error
	GetMetadata(ctx context.Context, userId int) ([]ImageMetadata, error)
	DeleteMetadata(ctx context.Context, key string) error
}

type s3Repository interface {
	PutObject(ctx context.Context, bucket, key string, object []byte) error
	GetObjectUrl(ctx context.Context, bucket, key string, limit time.Duration) (string, error)
	DeleteObject(ctx context.Context, bucket, key string) error
}

type ImageService struct {
	repo   imageRepository
	s3Repo s3Repository
}

func NewImageService(r imageRepository, s3Repo s3Repository) *ImageService {
	return &ImageService{repo: r, s3Repo: s3Repo}
}

func (s *ImageService) Upload(ctx context.Context, bucket, imageType string, object []byte) error {
	key := newS3Key(imageType)

	err := s.s3Repo.PutObject(ctx, bucket, key, object)
	if err != nil {
		return err
	}

	err = s.repo.AddMetadata(ctx, 1, imageType, key)
	if err != nil {
		return err
	}

	return nil
}

func (s *ImageService) GetImages(ctx context.Context, userId int) ([]ImageMetadata, error) {
	metadata, err := s.repo.GetMetadata(ctx, userId)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(metadata); i++ {
		s3url, err := s.s3Repo.GetObjectUrl(ctx, "images", metadata[i].Key, 80*time.Second)
		if err != nil {
			return nil, err
		}

		metadata[i].S3Url = s3url
	}

	return metadata, nil
}

func (s *ImageService) DeleteImage(ctx context.Context, bucket, key string) error {
	err := s.s3Repo.DeleteObject(ctx, bucket, key)
	if err != nil {
		return err
	}

	err = s.repo.DeleteMetadata(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
