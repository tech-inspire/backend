package consumer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	embeddingsv1 "github.com/tech-inspire/api-contracts/api/gen/go/embeddings/v1"
	"github.com/tech-inspire/backend/search-service/pkg/logger"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

type ImageEmbeddingsUpdatesConsumerProcessor interface {
	ProcessImageEmbeddingsUpdate(ctx context.Context, postID uuid.UUID, imageEmbeddings []float32) error
}
type ImageEmbeddingsUpdatesConsumer struct {
	js nats.JetStreamContext
}

func NewImageEmbeddingsUpdatesConsumer(js nats.JetStreamContext) (ImageEmbeddingsUpdatesConsumer, error) {
	return ImageEmbeddingsUpdatesConsumer{
		js: js,
	}, nil
}

func (c ImageEmbeddingsUpdatesConsumer) Start(lc fx.Lifecycle, processor ImageEmbeddingsUpdatesConsumerProcessor) error {
	process := func(msg *nats.Msg) error {
		var event embeddingsv1.PostEmbeddingsUpdatedEvent
		if err := proto.Unmarshal(msg.Data, &event); err != nil {
			return fmt.Errorf("unmarshal post created event: %w", err)
		}

		postID, err := uuid.Parse(event.PostId)
		if err != nil {
			return fmt.Errorf("parse post id: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err = processor.ProcessImageEmbeddingsUpdate(ctx, postID, event.EmbeddingVector)
		if err != nil {
			return fmt.Errorf("handle event: %w", err)
		}

		err = msg.Ack()
		if err != nil {
			return fmt.Errorf("ack event: %w", err)
		}

		return nil
	}

	shutDownCtx, cancel := context.WithCancel(context.Background())

	sub, err := c.js.QueueSubscribe(
		"posts.*.embeddings_updated",
		"posts-service-embeddings-workers",
		func(msg *nats.Msg) {
			slog.Debug("received embeddings update event", slog.String("subject", msg.Subject))
			if err := process(msg); err != nil {
				slog.Error("failed to process embeddings updated event",
					slog.String("subject", msg.Subject),
					logger.Error(err),
				)
			}
		},
		nats.Durable("posts-service-consumer-embeddings"),
		nats.ManualAck(),
		nats.Context(shutDownCtx),
	)
	if err != nil {
		cancel()
		return fmt.Errorf("subscribe: %w", err)
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			cancel()

			err = sub.Drain()
			if err != nil {
				return fmt.Errorf("drain subscription: %w", err)
			}

			return nil
		},
	})

	return nil
}
