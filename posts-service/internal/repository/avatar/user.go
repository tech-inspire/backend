package avatarstorage

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
	"github.com/tech-inspire/backend/posts-service/internal/service/dto"
)

func (AvatarStorage) userAvatarObjectName(userID uuid.UUID) string {
	return fmt.Sprintf("avatars/user_%s", userID)
}

func (fs AvatarStorage) UploadUserAvatar(ctx context.Context, params dto.UploadUserAvatar) (path string, err error) {
	objectName := fs.userAvatarObjectName(params.UserID)

	_, err = fs.client.PutObject(ctx, &s3.PutObjectInput{
		Body:          bytes.NewReader(params.Data),
		Bucket:        &fs.bucketName,
		Key:           &objectName,
		ACL:           types.ObjectCannedACLPublicRead,
		ContentLength: &params.ImageSize,
		ContentType:   &params.ContentType,
	})
	if err != nil {
		return "", fmt.Errorf("put object %s: %w", objectName, err)
	}

	return objectName, nil
}

func (fs AvatarStorage) GetUserAvatarURL(ctx context.Context, userID uuid.UUID) (string, error) {
	objectName := fs.userAvatarObjectName(userID)

	_, err := fs.client.GetObjectAttributes(ctx, &s3.GetObjectAttributesInput{
		Bucket: &fs.bucketName,
		Key:    &objectName,
	})
	if err != nil {
		return "", fmt.Errorf("get presigned url for %s: %w", objectName, err)
	}

	objectURL, err := url.JoinPath(fs.endpoint, fs.bucketName, objectName)
	if err != nil {
		return "", fmt.Errorf("get url for object: %w", err)
	}

	return objectURL, nil
}

func (fs AvatarStorage) DeleteUserAvatar(ctx context.Context, userID uuid.UUID) error {
	objectName := fs.userAvatarObjectName(userID)

	_, err := fs.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &fs.bucketName,
		Key:    &objectName,
	})
	if err != nil {
		return fmt.Errorf("remove object %s: %w", objectName, err)
	}

	return nil
}
