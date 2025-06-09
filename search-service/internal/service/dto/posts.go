package dto

import (
	"time"

	"github.com/google/uuid"
)

type PostCreatedEvent struct {
	PostID   uuid.UUID
	AuthorID uuid.UUID

	ImageKey    string
	ImageWidth  uint32
	ImageHeight uint32

	Description string
	CreatedAt   time.Time
}

type Iterator interface {
	Close() error
	Err() error
	PostIDs(yield func(uuid.UUID) bool)
}
