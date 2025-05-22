package scylla

import (
	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/tech-inspire/backend/posts-service/internal/models"
	"github.com/tech-inspire/backend/posts-service/pkg/generics"
)

func (p ImageVariant) toModel() models.ImageVariant {
	return models.ImageVariant{
		VariantType: p.VariantType,
		URL:         p.URL,
		Width:       p.Width,
		Height:      p.Height,
		Size:        p.Size,
	}
}

func (p *Post) toModel() *models.Post {
	return &models.Post{
		PostID:                   uuid.UUID(p.PostID),
		AuthorID:                 uuid.UUID(p.AuthorID),
		Images:                   generics.Convert(p.Images, ImageVariant.toModel),
		SoundCloudSongURL:        p.SoundCloudSong,
		SoundCloudSongStartMilli: p.SoundCloudSongStart,
		Description:              p.Description,
		CreatedAt:                p.CreatedAt,
	}
}

func imageVariantFromModel(p models.ImageVariant) ImageVariant {
	return ImageVariant{
		VariantType: p.VariantType,
		URL:         p.URL,
		Width:       p.Width,
		Height:      p.Height,
		Size:        p.Size,
	}
}

func postFromModel(p *models.Post) *Post {
	return &Post{
		PostID:              gocql.UUID(p.PostID),
		AuthorID:            gocql.UUID(p.AuthorID),
		Images:              generics.Convert(p.Images, imageVariantFromModel),
		SoundCloudSong:      p.SoundCloudSongURL,
		SoundCloudSongStart: p.SoundCloudSongStartMilli,
		Description:         p.Description,
		CreatedAt:           p.CreatedAt,
	}
}
