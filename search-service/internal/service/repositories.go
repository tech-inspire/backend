package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/tech-inspire/backend/search-service/internal/service/dto"
)

type TextEmbeddingsGenerator interface {
	GenerateTextEmbeddings(ctx context.Context, text string) ([]float32, error)
}

type ImageEmbeddingsTaskManager interface {
	SendGenerateImageEmbeddingsTask(ctx context.Context, postID uuid.UUID, imageURL string) error
}

type SearchRepository interface {
	SearchPosts(ctx context.Context, input dto.ProcessedSearchPostsParams) ([]dto.SearchResult, error)
	UpsertPost(ctx context.Context, params dto.CreatePostParams) error
	UpsertImageEmbeddings(ctx context.Context, postID uuid.UUID, embeddings []float32) error
	DeletePostInfo(ctx context.Context, postID uuid.UUID) error
}
