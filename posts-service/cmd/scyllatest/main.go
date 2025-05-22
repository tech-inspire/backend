package main

import (
	"context"
	"log"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/scylladb/gocqlx/v3"
	"github.com/tech-inspire/backend/posts-service/internal/models"
	"github.com/tech-inspire/backend/posts-service/internal/repository/scylla"
	"github.com/tech-inspire/backend/posts-service/migrations"
)

func main() {
	cluster := gocql.NewCluster(
		"127.0.0.1:19042",
		"127.0.0.1:19043",
		"127.0.0.1:19044",
		"127.0.0.1:19045",
		"127.0.0.1:19046",
	)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username:              "cassandra",
		Password:              "cassandra",
		AllowedAuthenticators: nil,
	}
	cluster.Keyspace = "posts"

	session, err := gocqlx.WrapSession(cluster.CreateSession())
	if err != nil {
		log.Fatal(err)
	}

	repo := scylla.NewPostsRepository(session)

	ctx := context.Background()

	err = migrations.ApplyMigrations(ctx, session)
	if err != nil {
		log.Fatal(err)
	}

	postID := uuid.New()

	err = repo.Create(ctx, &models.Post{
		PostID:   postID,
		AuthorID: uuid.New(),
		Images: []models.ImageVariant{
			{
				VariantType: "orig",
				URL:         "image",
				Width:       1200,
				Height:      700,
				Size:        12345,
			},
		},
		SoundCloudSongURL:        nil,
		SoundCloudSongStartMilli: nil,
		Description:              "my photo",
		CreatedAt:                time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

	p, err := repo.GetByID(ctx, postID)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(p)
}
