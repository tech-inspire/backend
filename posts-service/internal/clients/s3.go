package clients

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tech-inspire/backend/posts-service/internal/config"
	"go.uber.org/fx"
)

func NewS3Client(lc fx.Lifecycle, config *config.Config) (*s3.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithBaseEndpoint(config.S3.Endpoint),
	)
	if err != nil {
		return nil, err
	}

	// Create an Amazon S3 service client
	client := s3.NewFromConfig(cfg)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
				Bucket: aws.String(config.S3.BucketName),
			})
			if err != nil {
				return fmt.Errorf("check bucket %s exists: %w", config.S3.BucketName, err)
			}

			return nil
		},
	})

	return client, nil
}
