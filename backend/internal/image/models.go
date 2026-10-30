package image

import "time"

type ImageMetadata struct {
	Id      int       `json:"id"`
	Type    string    `json:"type"`
	Key     string    `json:"key"`
	S3Url   string    `json:"s3_url"`
	AddedAt time.Time `json:"added_at"`
	UserId  int       `json:"user_id"`
}
