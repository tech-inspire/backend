package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	"github.com/tech-inspire/backend/search-service/internal/models"
	"github.com/tech-inspire/backend/search-service/internal/service/dto"
	"github.com/tech-inspire/backend/search-service/pkg/generics"
)

type SearchRepository struct {
	pool *pgxpool.Pool
}

func NewSearchRepository(pool *pgxpool.Pool) *SearchRepository {
	return &SearchRepository{pool: pool}
}

func (r SearchRepository) UpsertPost(ctx context.Context, params dto.CreatePostParams) error {
	sb := sqlbuilder.NewInsertBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.InsertInto("posts_search_info").
		Cols("post_id", "author_id", "description", "image_path", "image_width", "image_height").
		Values(params.PostID, params.AuthorID, params.Description, params.ImagePath, params.ImageWidth, params.ImageHeight)
	query, args := sb.Build()
	slog.Debug("generated insert query", slog.String("query", query), slog.Any("args", args))

	_, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}

	return nil
}

func (r SearchRepository) DeletePostInfo(ctx context.Context, postID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM posts_search_info WHERE post_id = $1", postID)
	if err != nil {
		return fmt.Errorf("delete post: %w", err)
	}

	return nil
}

func (r SearchRepository) UpsertImageEmbeddings(ctx context.Context, postID uuid.UUID, embeddings []float32) error {
	v := pgvector.NewVector(embeddings)

	_, err := r.pool.Exec(ctx, "UPDATE posts_search_info SET image_embedding = $1, updated_at = NOW() WHERE post_id = $2", v, postID)
	if err != nil {
		return fmt.Errorf("update image_embedding: %w", err)
	}

	return nil
}

func applySearchParams(sb *sqlbuilder.SelectBuilder, params dto.ProcessedSearchPostsParams) (conditions []string, similarityUsed bool) {
	if params.AuthorID != nil {
		conditions = append(conditions, sb.Equal("author_id", *params.AuthorID))
	}

	if len(params.TextEmbeddings) > 0 {
		v := pgvector.NewVector(params.TextEmbeddings)

		column := fmt.Sprintf("image_embedding <=> (%s)::vector", sb.Var(v))
		sb.SelectMore(fmt.Sprintf("%s AS similarity_score", column))
		similarityUsed = true
	}

	if params.ReferencePostID != nil {
		column := fmt.Sprintf(
			"image_embedding <=> (SELECT r.image_embedding FROM posts_search_info r WHERE r.post_id = %s)",
			sb.Var(*params.ReferencePostID),
		)
		sb.SelectMore(fmt.Sprintf("%s AS similarity_score", column))
		similarityUsed = true
	}

	if params.PhotoOrientation != nil {
		minRatio, maxRation, ok := models.GetOrientationRange(*params.PhotoOrientation)
		if ok {
			conditions = append(conditions,
				sb.Between("ratio", minRatio, maxRation),
			)
		}
	}

	return conditions, similarityUsed
}

type searchPostsRow struct {
	PostID          uuid.UUID `db:"post_id"`
	SimilarityScore *float32  `db:"similarity_score"`
}

func (r SearchRepository) SearchPosts(ctx context.Context, input dto.ProcessedSearchPostsParams) ([]dto.SearchResult, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.SetFlavor(sqlbuilder.PostgreSQL)
	sb.Select("post_id").From("posts_search_info")

	conditions, similarityUsed := applySearchParams(sb, input)
	sb.Where(conditions...)

	if similarityUsed {
		sb.OrderBy("similarity_score ASC")
	}
	sb.OrderBy(fmt.Sprintf("%s %s", input.SearchSort, input.SearchOrder))

	sb.Offset(int(input.Offset))
	sb.Limit(int(input.Limit))

	query, args := sb.Build()

	slog.Info("generated search query", slog.String("query", query), slog.Any("args", args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []dto.SearchResult{}, nil
		}
		return nil, fmt.Errorf("query: %w", err)
	}

	results, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[searchPostsRow])
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("collect query: %w", err)
	}

	return generics.Convert(results, func(row searchPostsRow) dto.SearchResult {
		return dto.SearchResult(row)
	}), nil
}
