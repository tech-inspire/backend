package handlers

import (
	"context"

	"github.com/tech-inspire/backend/search-service/internal/service/dto"
)

type SearchService interface {
	SearchImages(ctx context.Context, params dto.SearchPostsParams) ([]dto.SearchResult, error)
}
