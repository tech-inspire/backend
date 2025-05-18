package dto

import (
	"github.com/google/uuid"
	"github.com/tech-inspire/service/auth-service/internal/models"
)

type RegisterParams struct {
	Email    string
	Username string
	Name     string

	Password string
}

type RegisterOutput struct {
	ConfirmationRequired bool
	LoginOutput          *LoginOutput
}

type CreateUserParams struct {
	UserID       uuid.UUID
	Email        string
	Name         string
	Username     string
	PasswordHash []byte
	Description  string
}

type LoginOutput struct {
	User    *models.User
	Session *models.Session
}
