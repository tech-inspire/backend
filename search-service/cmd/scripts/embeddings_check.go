package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	pgxvec "github.com/pgvector/pgvector-go/pgx"
	"github.com/tech-inspire/backend/search-service/internal/clients"
	"github.com/tech-inspire/backend/search-service/internal/config"
	natsrepo "github.com/tech-inspire/backend/search-service/internal/repository/nats"
	"github.com/tech-inspire/backend/search-service/internal/repository/postgres"
)

type Config struct {
	config.Database
	config.Nats
	config.ImageEmbeddings
	PostsServiceURL string `env:"POSTS_SERVICE_URL,required"`
}

func main() {
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	connectCfg, err := pgxpool.ParseConfig(cfg.Database.PostgresDSN)
	if err != nil {
		log.Fatalf("parse postgres dsn: %w", err)
	}
	connectCfg.AfterConnect = pgxvec.RegisterTypes

	pool, err := pgxpool.NewWithConfig(context.Background(), connectCfg)
	if err != nil {
		log.Fatalf("create pool: %w", err)
	}
	defer pool.Close()

	nc, err := nats.Connect(cfg.Nats.URL,
		nats.Name("search-service-embeddings-check"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(5*time.Second),
	)
	if err != nil {
		log.Fatalf("connect to nats: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("get jetstream context: %w", err)
	}
	defer nc.Drain()

	dispatcher, err := natsrepo.NewImageEmbeddingsEventDispatcher(js,
		cfg.Nats.PostsStreamName,
		cfg.ImageEmbeddings.ImageURLBasePath,
	)
	if err != nil {
		log.Fatalf("create dispatcher: %w", err)
	}

	client, err := clients.NewPostsServiceClient(cfg.PostsServiceURL)
	if err != nil {
		log.Fatalf("create posts service: %w", err)
	}

	repo := postgres.NewSearchRepository(pool)

	processor := Processor{
		repo:       repo,
		client:     client,
		dispatcher: dispatcher,
	}

	ctx := context.Background()

	if err = processor.process(ctx); err != nil {
		log.Fatalf("process posts: %w", err)
	}
}

type Processor struct {
	repo       *postgres.SearchRepository
	client     *clients.PostsServiceClient
	dispatcher *natsrepo.ImageEmbeddingsEventDispatcher
}

func (p Processor) process(ctx context.Context) error {
	iterator, err := p.repo.GetPostsWithoutImageEmbeddings(ctx)
	if err != nil {
		return fmt.Errorf("get posts iterator: %w", err)
	}

	total := 0
	defer func() {
		log.Printf("total processed posts: %d", total)
	}()

	for postID := range iterator.PostIDs {
		post, err := p.client.GetPostByID(ctx, postID)
		if err != nil {
			return fmt.Errorf("get post %s: %w", postID, err)
		}

		event, err := clients.PostCreatedEventFromPost(post)
		if err != nil {
			return fmt.Errorf("create post event: %w", err)
		}

		err = p.dispatcher.SendGenerateImageEmbeddingsTask(ctx, event.PostID, event.ImageKey)
		if err != nil {
			return fmt.Errorf("send image embeddings task: %w", err)
		}

		log.Println("Sent task for post: ", event.PostID)
		total++
	}

	if err = iterator.Err(); err != nil {
		return fmt.Errorf("iterator error: %w", err)
	}

	return nil
}
