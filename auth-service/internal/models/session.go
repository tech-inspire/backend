package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID uuid.UUID

	UserID    uuid.UUID
	Token     string
	CreatedAt time.Time
	ExpiresAt time.Time
}
