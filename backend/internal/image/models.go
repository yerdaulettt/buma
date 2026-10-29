package image

import "time"

type ImageMetadata struct {
	Id      int
	AddedAt time.Time
	Key     string
	UserId  int
	S3Url   string
}
