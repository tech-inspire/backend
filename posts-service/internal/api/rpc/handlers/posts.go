package handlers

import (
	"context"

	"connectrpc.com/connect"
	"github.com/tech-inspire/api-contracts/api/gen/go/posts/v1"
)

type PostsHandler struct{}

func (p PostsHandler) AddPost(ctx context.Context, c *connect.Request[postsv1.AddPostRequest]) (*connect.Response[postsv1.AddPostResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (p PostsHandler) GetPostByID(ctx context.Context, c *connect.Request[postsv1.GetPostByIDRequest]) (*connect.Response[postsv1.GetPostByIDResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (p PostsHandler) GetPosts(ctx context.Context, c *connect.Request[postsv1.GetPostsRequest]) (*connect.Response[postsv1.GetPostsResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (p PostsHandler) DeletePost(ctx context.Context, c *connect.Request[postsv1.DeletePostRequest]) (*connect.Response[postsv1.DeletePostResponse], error) {
	// TODO implement me
	panic("implement me")
}

func (p PostsHandler) GetUploadUrl(ctx context.Context, c *connect.Request[postsv1.GetUploadUrlRequest]) (*connect.Response[postsv1.GetUploadUrlResponse], error) {
	// TODO implement me
	panic("implement me")
}
