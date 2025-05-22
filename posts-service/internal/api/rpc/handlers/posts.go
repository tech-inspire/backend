package handlers

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/tech-inspire/api-contracts/api/gen/go/posts/v1"
	authmiddleware "github.com/tech-inspire/backend/auth-service/pkg/jwt/middleware"
	"github.com/tech-inspire/backend/posts-service/internal/proto"
	"github.com/tech-inspire/backend/posts-service/internal/service/dto"
	"github.com/tech-inspire/backend/posts-service/pkg/generics"
)

type PostsHandler struct {
	service PostsService
}

func NewPostsHandler(service PostsService) *PostsHandler {
	return &PostsHandler{service: service}
}

func (p PostsHandler) AddPost(ctx context.Context, c *connect.Request[postsv1.AddPostRequest]) (*connect.Response[postsv1.AddPostResponse], error) {
	userID := authmiddleware.GetUserInfo(ctx).UserID

	var soundcloudSongStart *int
	if c.Msg.SoundcloudSongStart != nil {
		start := int(*c.Msg.SoundcloudSongStart)
		soundcloudSongStart = &start
	}

	post, err := p.service.CreatePost(ctx, userID, dto.CreatePostParams{
		UploadSessionKey:         c.Msg.UploadSessionKey,
		ImageWidth:               int(c.Msg.ImageWidth),
		ImageHeight:              int(c.Msg.ImageHeight),
		ImageSize:                c.Msg.ImageSize,
		SoundCloudSongURL:        c.Msg.SoundcloudSong,
		SoundCloudSongStartMilli: soundcloudSongStart,
		Description:              c.Msg.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	return connect.NewResponse(&postsv1.AddPostResponse{
		Post: proto.Post(post),
	}), nil
}

func (p PostsHandler) GetPostByID(ctx context.Context, c *connect.Request[postsv1.GetPostByIDRequest]) (*connect.Response[postsv1.GetPostByIDResponse], error) {
	postID, err := uuid.Parse(c.Msg.PostId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse post_id: %w", err))
	}

	post, err := p.service.GetPostByID(ctx, postID)
	if err != nil {
		return nil, fmt.Errorf("get post %s: %w", postID, err)
	}

	return connect.NewResponse(&postsv1.GetPostByIDResponse{
		Post: proto.Post(post),
	}), nil
}

func (p PostsHandler) GetPosts(ctx context.Context, c *connect.Request[postsv1.GetPostsRequest]) (*connect.Response[postsv1.GetPostsResponse], error) {
	var err error

	postIDs := make([]uuid.UUID, len(c.Msg.PostIds))
	for i := range c.Msg.PostIds {
		postIDs[i], err = uuid.Parse(c.Msg.PostIds[i])
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse id %s: %w", c.Msg.PostIds[i], err))
		}
	}

	posts, err := p.service.GetPostsByIDs(ctx, postIDs)
	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}

	return connect.NewResponse(&postsv1.GetPostsResponse{
		Posts: generics.Convert(posts, proto.Post),
	}), nil
}

func (p PostsHandler) DeletePost(ctx context.Context, c *connect.Request[postsv1.DeletePostRequest]) (*connect.Response[postsv1.DeletePostResponse], error) {
	userID := authmiddleware.GetUserInfo(ctx).UserID

	postID, err := uuid.Parse(c.Msg.PostId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse post_id: %w", err))
	}

	err = p.service.DeletePostByID(ctx, userID, postID)
	if err != nil {
		return nil, fmt.Errorf("delete post %s: %w", postID, err)
	}

	return &connect.Response[postsv1.DeletePostResponse]{}, nil
}

func (p PostsHandler) GetUploadUrl(ctx context.Context, c *connect.Request[postsv1.GetUploadUrlRequest]) (*connect.Response[postsv1.GetUploadUrlResponse], error) {
	userID := authmiddleware.GetUserInfo(ctx).UserID

	res, err := p.service.GenerateTempImageUpload(ctx, dto.GenerateImageUploadURLParams{
		UserID:      userID,
		ImageSize:   c.Msg.FileSize,
		ContentType: c.Msg.MimeType,
	})
	if err != nil {
		return nil, fmt.Errorf("generate temp image upload url: %w", err)
	}

	headers := make(map[string]string, len(res.PresignedURL.Headers))
	for k := range res.PresignedURL.Headers {
		headers[k] = res.PresignedURL.Headers.Get(k)
	}

	return connect.NewResponse(&postsv1.GetUploadUrlResponse{
		UploadUrl:        res.PresignedURL.URL,
		UploadSessionKey: res.Key,
		Method:           res.PresignedURL.Method,
		Headers:          headers,
	}), nil
}
