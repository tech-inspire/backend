package clients

import (
	"context"
	"fmt"
	"net/http"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	postsv1 "github.com/tech-inspire/api-contracts/api/gen/go/posts/v1"
	"github.com/tech-inspire/api-contracts/api/gen/go/posts/v1/postsv1connect"
	"github.com/tech-inspire/backend/search-service/internal/service/dto"
)

type PostsServiceClient struct {
	client postsv1connect.PostsServiceClient
}

func NewPostsServiceClient(url string) (*PostsServiceClient, error) {
	c := postsv1connect.NewPostsServiceClient(
		http.DefaultClient,
		url,
	)

	client := &PostsServiceClient{
		client: c,
	}

	return client, nil
}

func (e PostsServiceClient) GetPostByID(ctx context.Context, postID uuid.UUID) (*postsv1.Post, error) {
	resp, err := e.client.GetPostByID(ctx, connect.NewRequest(&postsv1.GetPostByIDRequest{
		PostId: postID.String(),
	}))
	if err != nil {
		return nil, fmt.Errorf("posts service: GetPostByID(%s): %s", postID, err)
	}

	return resp.Msg.Post, nil
}

func PostCreatedEventFromPost(post *postsv1.Post) (dto.PostCreatedEvent, error) {
	postID, err := uuid.Parse(post.PostId)
	if err != nil {
		return dto.PostCreatedEvent{}, fmt.Errorf("parse post id: %w", err)
	}

	authorID, err := uuid.Parse(post.AuthorId)
	if err != nil {
		return dto.PostCreatedEvent{}, fmt.Errorf("parse post id: %w", err)
	}

	if len(post.Images) == 0 {
		return dto.PostCreatedEvent{}, fmt.Errorf("post images is empty")
	}

	image := post.Images[0]

	return dto.PostCreatedEvent{
		PostID:      postID,
		AuthorID:    authorID,
		ImageKey:    image.Url,
		ImageWidth:  uint32(image.Width),
		ImageHeight: uint32(image.Height),
		Description: post.Description,
		CreatedAt:   post.CreatedAt.AsTime(),
	}, nil
}
