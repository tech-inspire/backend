package clients

import (
	"context"
	"fmt"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	transport "github.com/aws/smithy-go/endpoints"
	"github.com/tech-inspire/backend/auth-service/internal/config"
	"go.uber.org/fx"
)

func NewS3Client(lc fx.Lifecycle, config *config.Config) (*s3.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithBaseEndpoint(config.S3.Endpoint),
		awsconfig.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if config.S3.MinioResolveMode {
			o.EndpointResolverV2 = &MinioBucketResolver{config.S3.Endpoint}
		}
	})

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

type MinioBucketResolver struct {
	URL string
}

func (r *MinioBucketResolver) ResolveEndpoint(_ context.Context, params s3.EndpointParameters) (transport.Endpoint, error) {
	u, err := url.Parse(r.URL)
	if err != nil {
		return transport.Endpoint{}, fmt.Errorf("parse S3 URL: %w", err)
	}

	u = u.JoinPath(*params.Bucket)
	return transport.Endpoint{URI: *u}, nil
}
