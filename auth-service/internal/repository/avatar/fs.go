package avatarstorage

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tech-inspire/service/auth-service/internal/config"
)

type AvatarStorage struct {
	client        *s3.Client
	presignClient *s3.PresignClient
	endpoint      string
	bucketName    string
}

func New(cfg *config.Config, client *s3.Client) (*AvatarStorage, error) {
	return &AvatarStorage{
		client:        client,
		presignClient: s3.NewPresignClient(client),
		bucketName:    cfg.S3.BucketName,
	}, nil
}
