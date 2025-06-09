package nats

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	embeddingsv1 "github.com/tech-inspire/api-contracts/api/gen/go/embeddings/v1"
	"google.golang.org/protobuf/proto"
)

type ImageEmbeddingsEventDispatcher struct {
	js         nats.JetStreamContext
	streamName string

	postImageBasePath string
}

func NewImageEmbeddingsEventDispatcher(js nats.JetStreamContext, streamName string, postImageBasePath string) (*ImageEmbeddingsEventDispatcher, error) {
	return &ImageEmbeddingsEventDispatcher{
		js:                js,
		streamName:        streamName,
		postImageBasePath: postImageBasePath,
	}, nil
}

func (d *ImageEmbeddingsEventDispatcher) SendGenerateImageEmbeddingsTask(ctx context.Context, postID uuid.UUID, imageURL string) error {
	fullImageURL, err := url.JoinPath(d.postImageBasePath, imageURL)
	if err != nil {
		return fmt.Errorf("create image url: %w", err)
	}

	msg := &embeddingsv1.GeneratePostEmbeddingsEvent{
		PostId:   postID.String(),
		ImageUrl: fullImageURL,
	}

	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	subject := fmt.Sprintf("posts.%s.generate_embeddings", postID)

	pubOpts := []nats.PubOpt{
		nats.Context(ctx),
		nats.ExpectStream(d.streamName),
	}
	if _, err = d.js.Publish(subject, data, pubOpts...); err != nil {
		return fmt.Errorf("publish %s: %w", subject, err)
	}

	return nil
}
