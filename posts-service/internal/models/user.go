package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID

	Name  string
	Email string

	Username    string
	Description string

	IsAdmin bool

	CreatedAt time.Time

	AvatarURL *string
}
