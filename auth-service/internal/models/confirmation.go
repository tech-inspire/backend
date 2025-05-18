package models

import (
	"time"

	"github.com/google/uuid"
)

type ConfirmationUserData struct {
	ConfirmationCode string

	Email        string
	Username     string
	Name         string
	PasswordHash string
	ExpiresAt    time.Time // To track expiration time for the confirmation code
}

type ResetPasswordData struct {
	UserID    uuid.UUID
	Code      string
	Email     string
	ExpiresAt time.Time
}
