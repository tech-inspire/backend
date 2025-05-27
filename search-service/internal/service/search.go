package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tech-inspire/backend/search-service/internal/service/dto"
)

type SearchService struct {
	repo        SearchRepository
	embeddings  TextEmbeddingsGenerator
	taskManager ImageEmbeddingsTaskManager
}

func NewSearchService(repo SearchRepository, embeddings TextEmbeddingsGenerator, taskManager ImageEmbeddingsTaskManager) *SearchService {
	return &SearchService{repo: repo, embeddings: embeddings, taskManager: taskManager}
}

func (s *SearchService) ProcessEventUpdated(ctx context.Context, event dto.PostCreatedEvent) error {
	err := s.repo.UpsertPost(ctx, dto.CreatePostParams{
		PostID:      event.PostID,
		AuthorID:    event.AuthorID,
		Description: event.Description,
		ImagePath:   event.ImageKey,
		ImageWidth:  event.ImageWidth,
		ImageHeight: event.ImageHeight,
	})
	if err != nil {
		return fmt.Errorf("update post info: %w", err)
	}

	err = s.taskManager.SendGenerateImageEmbeddingsTask(ctx, event.PostID, event.ImageKey)
	if err != nil {
		return fmt.Errorf("send image embeddings task: %w", err)
	}

	return nil
}

func (s *SearchService) ProcessImageEmbeddingsUpdate(ctx context.Context, postID uuid.UUID, imageEmbeddings []float32) error {
	return s.repo.UpsertImageEmbeddings(ctx, postID, imageEmbeddings)
}

func (s *SearchService) SearchImages(ctx context.Context, params dto.SearchPostsParams) ([]dto.SearchResult, error) {
	processedParams := dto.ProcessedSearchPostsParams{
		TextEmbeddings: nil,
		SearchParams:   params.SearchParams,
	}

	if params.TextQuery != nil {
		embeddings, err := s.embeddings.GenerateTextEmbeddings(ctx, *params.TextQuery)
		if err != nil {
			return nil, fmt.Errorf("generate text embeddings: %w", err)
		}

		processedParams.TextEmbeddings = embeddings
	}

	return s.repo.SearchPosts(ctx, processedParams)
}

func (s *SearchService) ProcessEventDeleted(ctx context.Context, postID uuid.UUID) error {
	return s.repo.DeletePostInfo(ctx, postID)
}
