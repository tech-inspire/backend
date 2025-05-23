package clients

import (
	"context"
	"fmt"

	embeddingsv1 "github.com/tech-inspire/api-contracts/api/gen/go/embeddings/v1"
	"github.com/tech-inspire/backend/search-service/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EmbeddingServiceClient struct {
	client embeddingsv1.EmbeddingsServiceClient
}

func NewEmbeddingServiceClient(cfg *config.Config) (*EmbeddingServiceClient, error) {
	conn, err := grpc.NewClient(cfg.EmbeddingsClient.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("grpc: create client: %w", err)
	}

	c := embeddingsv1.NewEmbeddingsServiceClient(conn)

	return &EmbeddingServiceClient{
		client: c,
	}, nil
}

func (e EmbeddingServiceClient) GenerateTextEmbeddings(ctx context.Context, text string) ([]float32, error) {
	resp, err := e.client.GenerateTextEmbeddings(ctx, &embeddingsv1.GenerateTextEmbeddingsRequest{
		Text: text,
	})
	if err != nil {
		return nil, fmt.Errorf("embeddings service: GenerateTextEmbeddings: %s", err)
	}

	return resp.EmbeddingVector, nil
}
