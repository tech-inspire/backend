package handlers

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	searchv1 "github.com/tech-inspire/api-contracts/api/gen/go/search/v1"
	"github.com/tech-inspire/backend/search-service/internal/service/dto"
	"github.com/tech-inspire/backend/search-service/pkg/generics"
)

type SearchHandler struct {
	SearchService SearchService
}

func NewSearchHandler(searchService SearchService) *SearchHandler {
	return &SearchHandler{SearchService: searchService}
}

func (h SearchHandler) SearchPosts(ctx context.Context, req *connect.Request[searchv1.SearchImagesRequest]) (*connect.Response[searchv1.SearchImagesResponse], error) {
	params := dto.SearchPostsParams{
		SearchParams: dto.SearchParams{
			SearchOrder: dto.Desc,
			SearchSort:  dto.CreatedAt,
			Offset:      req.Msg.Offset,
			Limit:       req.Msg.Limit,
		},
	}

	if req.Msg.AuthorId != nil {
		authorID, err := uuid.Parse(*req.Msg.AuthorId)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse author_id: %w", err))
		}

		params.AuthorID = &authorID
	}

	switch v := req.Msg.SearchBy.(type) {
	case *searchv1.SearchImagesRequest_TextQuery:
		params.TextQuery = &v.TextQuery
	case *searchv1.SearchImagesRequest_ReferencePostId:
		postID, err := uuid.Parse(v.ReferencePostId)
		if err != nil {
			return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("parse post_id: %w", err))
		}

		params.ReferencePostID = &postID
	}

	results, err := h.SearchService.SearchImages(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("search images: %w", err)
	}

	return connect.NewResponse(&searchv1.SearchImagesResponse{
		Results: generics.Convert(results, func(res dto.SearchResult) *searchv1.SearchResult {
			return &searchv1.SearchResult{
				PostId:     res.PostID.String(),
				Similarity: res.SimilarityScore,
			}
		}),
		Limit:  req.Msg.Limit,
		Offset: req.Msg.Offset,
	}), nil
}
