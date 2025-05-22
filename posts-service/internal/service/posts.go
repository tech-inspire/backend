package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/tech-inspire/backend/posts-service/internal/apperrors"
	"github.com/tech-inspire/backend/posts-service/internal/models"
	"github.com/tech-inspire/backend/posts-service/internal/service/dto"
)

type PostsService struct {
	repo          PostsRepository
	imageStorage  ImageStorage
	pendingImages PendingImagesRepository
	dispatcher    PostsEventDispatcher
}

func NewPostsService(
	repo PostsRepository,
	imageStorage ImageStorage,
	pendingImages PendingImagesRepository,
	dispatcher PostsEventDispatcher,
) *PostsService {
	return &PostsService{
		repo:          repo,
		imageStorage:  imageStorage,
		pendingImages: pendingImages,
		dispatcher:    dispatcher,
	}
}

func (p PostsService) GenerateTempImageUpload(ctx context.Context, params dto.GenerateImageUploadURLParams) (*dto.GeneratedImageUpload, error) {
	const expireTime = time.Minute * 15
	expires := time.Now().Add(expireTime)

	res, err := p.imageStorage.GenerateTempImageUpload(ctx, params, expireTime)
	if err != nil {
		return nil, fmt.Errorf("generate temp image upload: %w", err)
	}

	err = p.pendingImages.Add(ctx, res.Key, expires)
	if err != nil {
		return nil, fmt.Errorf("add images to pending list: %w", err)
	}

	return res, nil
}

func (p PostsService) UpdatePostByID(ctx context.Context, userID uuid.UUID, postID uuid.UUID, params dto.UpdatePostParams) error {
	return errors.New("not implemented")
}

func (p PostsService) CreatePost(ctx context.Context, userID uuid.UUID, params dto.CreatePostParams) (*models.Post, error) {
	err := p.pendingImages.Remove(ctx, params.UploadSessionKey)
	if err != nil {
		return nil, fmt.Errorf("remove image from pending list: %w", err)
	}

	postImage, err := p.imageStorage.CreatePostImage(ctx, params.UploadSessionKey, userID)
	if err != nil {
		return nil, fmt.Errorf("image storage: create post image: %w", err)
	}

	post := &models.Post{
		PostID:   uuid.Must(uuid.NewV7()),
		AuthorID: userID,
		Images: []models.ImageVariant{
			{
				VariantType: models.Original,
				URL:         postImage.PostKey,
				Width:       params.ImageWidth,
				Height:      params.ImageHeight,
				Size:        params.ImageSize,
			},
		},
		SoundCloudSongURL:        params.SoundCloudSongURL,
		SoundCloudSongStartMilli: params.SoundCloudSongStartMilli,
		Description:              params.Description,
		CreatedAt:                time.Now(),
	}

	err = p.repo.CreatePost(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("create post: %w", err)
	}

	err = p.dispatcher.DispatchPostCreatedEvent(ctx, post)
	if err != nil {
		return nil, fmt.Errorf("dispatch post created event: %w", err)
	}

	return post, nil
}

func (p PostsService) GetPostByID(ctx context.Context, postID uuid.UUID) (*models.Post, error) {
	return p.repo.GetPostByID(ctx, postID)
}

func (p PostsService) GetPostsByIDs(ctx context.Context, postIDs []uuid.UUID) ([]*models.Post, error) {
	return p.repo.GetPostsByIDs(ctx, postIDs)
}

func (p PostsService) DeletePostByID(ctx context.Context, userID uuid.UUID, postID uuid.UUID) error {
	post, err := p.repo.GetPostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("get post: %w", err)
	}

	if post.AuthorID != userID {
		return apperrors.ErrForbidden
	}

	err = p.repo.DeletePostByID(ctx, postID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	return nil
}
