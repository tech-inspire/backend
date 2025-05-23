package models

import (
	"time"

	"github.com/google/uuid"
)

type VariantType string

const (
	Original  VariantType = "original"
	Thumbnail VariantType = "thumbnail"
)

type ImageVariant struct {
	VariantType VariantType
	URL         string
	Width       int
	Height      int
	Size        int32
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
