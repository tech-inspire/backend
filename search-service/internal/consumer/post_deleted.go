package consumer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	postsv1 "github.com/tech-inspire/api-contracts/api/gen/go/posts/v1"
	"github.com/tech-inspire/backend/search-service/pkg/logger"
	"go.uber.org/fx"
	"google.golang.org/protobuf/proto"
)

func StartPostDeletedEventsConsumer(js nats.JetStreamContext, lc fx.Lifecycle, processor PostsEventProcessor) error {
	process := func(msg *nats.Msg) error {
		var event postsv1.PostDeletedEvent
		if err := proto.Unmarshal(msg.Data, &event); err != nil {
			return fmt.Errorf("unmarshal post deleted event: %w", err)
		}

		postID, err := uuid.Parse(event.Post.PostId)
		if err != nil {
			return fmt.Errorf("parse post id: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		err = processor.ProcessEventDeleted(ctx, postID)
		if err != nil {
			return fmt.Errorf("handle post deleted event: %w", err)
		}

		err = msg.Ack()
		if err != nil {
			return fmt.Errorf("ack event: %w", err)
		}

		slog.Info("processed posts deleted event", slog.String("sub", msg.Subject))

		return nil
	}

	shutDownCtx, cancel := context.WithCancel(context.Background())

	sub, err := js.QueueSubscribe(
		"posts.*.deleted",
		"posts-service-posts-workers",
		func(msg *nats.Msg) {
			if err := process(msg); err != nil {
				slog.Error("failed to process post deleted event",
					slog.String("subject", msg.Subject),
					logger.Error(err),
				)
			}
		},
		nats.Durable("posts-service-consumer-posts-deleted"),
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
