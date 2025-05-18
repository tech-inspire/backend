package dto

import (
	"github.com/google/uuid"
)

type UploadUserAvatar struct {
	Data        []byte
	UserID      uuid.UUID
	ImageSize   int64
	ContentType string
}
