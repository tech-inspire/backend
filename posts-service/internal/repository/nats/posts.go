package nats

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	postsv1 "github.com/tech-inspire/api-contracts/api/gen/go/posts/v1"
	"github.com/tech-inspire/backend/posts-service/internal/config"
	"github.com/tech-inspire/backend/posts-service/internal/models"
	postsproto "github.com/tech-inspire/backend/posts-service/internal/proto"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PostsEventDispatcher struct {
	js         nats.JetStreamContext
	streamName string
}

func NewPostsEventDispatcher(cfg *config.Config) (*PostsEventDispatcher, error) {
	nc, err := nats.Connect(cfg.Nats.URL,
		nats.Name("posts-service-event-dispatcher"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("get jetstream context: %w", err)
	}

	return &PostsEventDispatcher{
		js:         js,
		streamName: cfg.Nats.PostsStreamName,
	}, nil
}

func (d *PostsEventDispatcher) DispatchPostCreatedEvent(ctx context.Context, post *models.Post) error {
	msg := &postsv1.PostCreatedEvent{
		CreatedAt: timestamppb.New(post.CreatedAt),
		Post:      postsproto.Post(post),
	}

	return d.publishEvent(ctx, post.PostID, "created", msg)
}

func (d *PostsEventDispatcher) DispatchPostUpdatedEvent(ctx context.Context, post *models.Post, updatedAt time.Time) error {
	msg := &postsv1.PostUpdatedEvent{
		UpdatedAt: timestamppb.New(updatedAt),
		Post:      postsproto.Post(post),
	}

	return d.publishEvent(ctx, post.PostID, "updated", msg)
}

func (d *PostsEventDispatcher) DispatchPostDeletedEvent(ctx context.Context, post *models.Post, deletedAt time.Time) error {
	msg := &postsv1.PostDeletedEvent{
		DeletedAt: timestamppb.New(deletedAt),
		Post:      postsproto.Post(post),
	}

	return d.publishEvent(ctx, post.PostID, "deleted", msg)
}

func (d *PostsEventDispatcher) publishEvent(ctx context.Context, postID uuid.UUID, action string, message proto.Message) error {
	payload, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal proto: %w", err)
	}

	subject := fmt.Sprintf("posts.%s.%s", postID, action)

	pubOpts := []nats.PubOpt{
		nats.Context(ctx),
		nats.ExpectStream(d.streamName),
	}
	if _, err = d.js.Publish(subject, payload, pubOpts...); err != nil {
		return fmt.Errorf("publish %s: %w", subject, err)
	}

	return nil
}
