package consumer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	postsv1 "github.com/tech-inspire/api-contracts/api/gen/go/posts/v1"
	"github.com/tech-inspire/backend/search-service/internal/service/dto"
	"github.com/tech-inspire/backend/search-service/pkg/logger"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

type PostsEventProcessor interface {
	ProcessEventUpdated(ctx context.Context, event dto.PostCreatedEvent) error
	ProcessEventDeleted(ctx context.Context, postID uuid.UUID) error
}

func StartPostCreatedEventsConsumer(js nats.JetStreamContext, lc fx.Lifecycle, processor PostsEventProcessor) error {
	process := func(msg *nats.Msg) error {
		var event postsv1.PostCreatedEvent
		if err := proto.Unmarshal(msg.Data, &event); err != nil {
			return fmt.Errorf("unmarshal post created event: %w", err)
		}

		postID, err := uuid.Parse(event.Post.PostId)
		if err != nil {
			return fmt.Errorf("parse post id: %w", err)
		}

		authorID, err := uuid.Parse(event.Post.AuthorId)
		if err != nil {
			return fmt.Errorf("parse post id: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if len(event.Post.Images) == 0 {
			return fmt.Errorf("post images is empty")
		}

		image := event.Post.Images[0]

		err = processor.ProcessEventUpdated(ctx, dto.PostCreatedEvent{
			PostID:      postID,
			AuthorID:    authorID,
			ImageKey:    image.Url,
			ImageWidth:  uint32(image.Width),
			ImageHeight: uint32(image.Height),
			Description: event.Post.Description,
			CreatedAt:   event.Post.CreatedAt.AsTime(),
		})
		if err != nil {
			return fmt.Errorf("handle post created event: %w", err)
		}

		err = msg.Ack()
		if err != nil {
			return fmt.Errorf("ack event: %w", err)
		}

		slog.Info("processed posts created event", slog.String("sub", msg.Subject))

		return nil
	}

	shutDownCtx, cancel := context.WithCancel(context.Background())

	sub, err := js.QueueSubscribe(
		"posts.*.created",
		"posts-service-posts-workers",
		func(msg *nats.Msg) {
			if err := process(msg); err != nil {
				slog.Error("failed to process post created event",
					slog.String("subject", msg.Subject),
					logger.Error(err),
				)
			}
		},
		nats.Durable("posts-service-consumer-posts"),
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
