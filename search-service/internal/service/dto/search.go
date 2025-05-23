package dto

import (
	"github.com/google/uuid"
	"github.com/tech-inspire/backend/search-service/internal/models"
)

type SearchOrder string

const (
	Asc  SearchOrder = "asc"
	Desc SearchOrder = "desc"
)

type SearchSort string

const (
	CreatedAt SearchSort = "created_at"
)

type SearchPostsParams struct {
	TextQuery *string
	SearchParams
}

type SearchParams struct {
	ReferencePostID  *uuid.UUID
	AuthorID         *uuid.UUID
	PhotoOrientation *models.PhotoOrientation

	SearchOrder SearchOrder
	SearchSort  SearchSort

	Offset uint32
	Limit  uint32
}

type ProcessedSearchPostsParams struct {
	TextEmbeddings []float32

	SearchParams
}

type SearchResult struct {
	PostID          uuid.UUID
	SimilarityScore *float32
}

type CreatePostParams struct {
	PostID      uuid.UUID
	AuthorID    uuid.UUID
	Description string
	ImagePath   string
	ImageWidth  uint32
	ImageHeight uint32
}
