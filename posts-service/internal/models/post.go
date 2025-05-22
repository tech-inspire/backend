package models

import (
	"time"

	"github.com/google/uuid"
)

type ImageVariant struct {
	VariantType string
	URL         string
	Width       int
	Height      int
	Size        int64
}

// Post maps to the posts_by_id table.
type Post struct {
	PostID                   uuid.UUID
	AuthorID                 uuid.UUID
	Images                   []ImageVariant
	SoundCloudSongURL        *string
	SoundCloudSongStartMilli *int
	Description              string
	CreatedAt                time.Time
}
