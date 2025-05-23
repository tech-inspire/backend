package clients

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tech-inspire/backend/search-service/internal/config"
)

func NewNatsJetstreamClient(cfg *config.Config) (nats.JetStreamContext, error) {
	nc, err := nats.Connect(cfg.Nats.URL,
		nats.Name("search-service"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("get jetstream context: %w", err)
	}

	return js, nil
}
