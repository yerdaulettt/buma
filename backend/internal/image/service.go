package image

import "context"

type imageRepository interface {
}

type s3Repository interface {
	PutObject(ctx context.Context, bucket, key string, object []byte) error
	DeleteObject(ctx context.Context, bucket, key string) error
}

type ImageService struct {
	repo   imageRepository
	s3Repo s3Repository
}

func NewImageService(r imageRepository, s3Repo s3Repository) *ImageService {
	return &ImageService{repo: r, s3Repo: s3Repo}
}

func (s *ImageService) Upload(ctx context.Context, bucket, key string, object []byte) error {
	err := s.s3Repo.PutObject(ctx, bucket, key, object)
	if err != nil {
		return err
	}

	return nil
}

func (s *ImageService) DeleteImage(ctx context.Context, bucket, key string) error {
	err := s.s3Repo.DeleteObject(ctx, bucket, key)
	if err != nil {
		return err
	}

	return nil
}
