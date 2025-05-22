package dto

import (
	"net/http"

	"github.com/google/uuid"
)

type GenerateImageUploadURLParams struct {
	UserID uuid.UUID

	ImageSize   int64
	ContentType string
}

type GeneratedImageUpload struct {
	PresignedURL PresignedURL
	Key          string
}

type PresignedURL struct {
	Method  string
	URL     string
	Headers http.Header
}

type CreatedPostImage struct {
	PostKey string
}
