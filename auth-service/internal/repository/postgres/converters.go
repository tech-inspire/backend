package postgres

import (
	"github.com/tech-inspire/service/auth-service/internal/models"
	"github.com/tech-inspire/service/auth-service/internal/repository/postgres/sqlc"
)

func userToModel(user sqlc.User) *models.User {
	return &models.User{
		ID:          user.UserID,
		Name:        user.Name,
		Email:       user.Email,
		Username:    user.Username,
		Description: user.Description,
		IsAdmin:     user.IsAdmin,
		CreatedAt:   user.CreatedAt,
		AvatarURL:   user.AvatarUrl,
	}
}
