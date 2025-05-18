package redis

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type UserConfirmationData struct {
	ConfirmationCode string `json:"confirmation_code"`

	Email        string    `json:"email"`
	Username     string    `json:"username"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"password_hash"`
	ExpiresAt    time.Time `json:"expires_at"` // To track expiration time for the confirmation code
}

type ResetPasswordData struct {
	UserID    uuid.UUID `json:"user_id"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}
