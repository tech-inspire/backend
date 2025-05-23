package imagestorage

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/tech-inspire/backend/posts-service/internal/config"
	"github.com/tech-inspire/backend/posts-service/internal/service/dto"
)

type ImageStorage struct {
	client    *s3.Client
	presigner *s3.PresignClient

	endpoint   string
	bucketName string
}

func New(cfg *config.Config, client *s3.Client) (*ImageStorage, error) {
	return &ImageStorage{
		client:     client,
		presigner:  s3.NewPresignClient(client),
		bucketName: cfg.S3.BucketName,
	}, nil
}

func (ImageStorage) imageObjectName(postID uuid.UUID) string {
	return fmt.Sprintf("images/post_%s", postID)
}

func (ImageStorage) tempImageObjectName(userID uuid.UUID) string {
	return fmt.Sprintf("tmp/images/%d_%s", time.Now().UnixNano(), userID)
}

func (fs ImageStorage) GenerateTempImageUpload(ctx context.Context, params dto.GenerateImageUploadURLParams, expire time.Duration) (*dto.GeneratedImageUpload, error) {
	objectName := fs.tempImageObjectName(params.UserID)

	res, err := fs.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:        &fs.bucketName,
		Key:           &objectName,
		ACL:           types.ObjectCannedACLPublicRead,
		ContentLength: aws.Int64(int64(params.ImageSize)),
		ContentType:   &params.ContentType,
	}, s3.WithPresignExpires(expire))
	if err != nil {
		return nil, fmt.Errorf("presign: put object %s: %w", objectName, err)
	}

	return &dto.GeneratedImageUpload{
		PresignedURL: dto.PresignedURL{
			Method:  res.Method,
			URL:     res.URL,
			Headers: res.SignedHeader,
		},
		Key: objectName,
	}, nil
}

func (fs ImageStorage) CreatePostImage(ctx context.Context, tempImageKey string, postID uuid.UUID) (*dto.CreatedPostImage, error) {
	objectName := fs.imageObjectName(postID)

	copySource := filepath.Join(fs.bucketName, tempImageKey)

	_, err := fs.client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     &fs.bucketName,
		CopySource: &copySource,
		Key:        &objectName,
		ACL:        types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return nil, fmt.Errorf("copy object %s to %s: %w", tempImageKey, objectName, err)
	}

	_, err = fs.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &fs.bucketName,
		Key:    &tempImageKey,
	})
	if err != nil {
		return nil, fmt.Errorf("delete object %s: %w", tempImageKey, err)
	}

	return &dto.CreatedPostImage{
		PostKey: objectName,
	}, nil
}

func (fs ImageStorage) DeletePostImage(ctx context.Context, postID uuid.UUID) error {
	objectName := fs.imageObjectName(postID)

	_, err := fs.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &fs.bucketName,
		Key:    &objectName,
	})
	if err != nil {
		return fmt.Errorf("remove object %s: %w", objectName, err)
	}

	return nil
}
