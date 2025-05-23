package scylla

import (
	"context"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/scylladb/gocqlx/v3"
	"github.com/scylladb/gocqlx/v3/qb"
	"github.com/tech-inspire/backend/posts-service/internal/models"
	"github.com/tech-inspire/backend/posts-service/internal/service/dto"
	"github.com/tech-inspire/backend/posts-service/pkg/generics"
)

type PostsRepository struct {
	session gocqlx.Session
}

func NewPostsRepository(session gocqlx.Session) *PostsRepository {
	return &PostsRepository{session: session}
}

func (r *PostsRepository) Create(ctx context.Context, p *models.Post) error {
	schemaPost := postFromModel(p)

	q := r.session.Query(postTable.Insert()).
		WithContext(ctx).
		BindStruct(schemaPost)
	if err := q.Err(); err != nil {
		return fmt.Errorf("insert query: bind post: %w", err)
	}

	if err := q.ExecRelease(); err != nil {
		return fmt.Errorf("insert query: exec release: %w", err)
	}
	return nil
}

// GetByID retrieves a Post by its post_id.
func (r *PostsRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Post, error) {
	var p Post

	query := r.session.
		Query(postTable.Select()).
		WithContext(ctx).
		Consistency(gocql.One).
		Bind(gocql.UUID(id))
	if err := query.GetRelease(&p); err != nil {
		return nil, fmt.Errorf("query: get post: %s: %w", query, err)
	}

	return p.toModel(), nil
}

// GetMany fetches multiple Posts by their IDs.
func (r *PostsRepository) GetMany(ctx context.Context, ids []uuid.UUID) ([]*models.Post, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	// or N requests to the nodes
	// https://lostechies.com/ryansvihla/2014/09/22/cassandra-query-patterns-not-using-the-in-query-for-multiple-partitions/

	// Build a single IN query
	stmt, names := qb.Select(postMetadata.Name).
		Columns(postMetadata.Columns...).
		Where(qb.In("post_id")).
		ToCql()

	cqlIDs := generics.Convert(ids, func(id uuid.UUID) gocql.UUID {
		return gocql.UUID(id)
	})
	q := r.session.Query(stmt, names).Bind(cqlIDs).WithContext(ctx)

	var posts []*Post
	if err := q.SelectRelease(&posts); err != nil {
		return nil, fmt.Errorf("get many posts: %w", err)
	}

	return generics.Convert(posts, (*Post).toModel), nil
}

// Update replaces an existing Post. It must include the post_id.
func (r *PostsRepository) Update(ctx context.Context, postID uuid.UUID, params dto.UpdatePostParams) (*models.Post, error) {
	// stmt, names := postTable.Update().Where(qb.Eq("post_id")).ToCql()
	// q := gocqlx.Query(r.session.Query(stmt), names).WithContext(ctx)
	// return q.BindStruct(p).ExecRelease()
	panic("implement me!")
}

// Delete removes a Post by its post_id.
func (r *PostsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// stmt, names := postTable.Delete().Where(qb.Eq("post_id")).ToCql()
	// q := gocqlx.Query(r.session.Query(stmt, id), names).WithContext(ctx)
	// return q.ExecRelease()

	panic("implement me!")
}
