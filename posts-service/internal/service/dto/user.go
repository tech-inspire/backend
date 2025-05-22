package dto

import (
	"github.com/tech-inspire/backend/posts-service/internal/models"
)

type GetUserByIDOutput struct {
	User *models.User
}

type GetCurrentUser struct {
	GetUserByIDOutput
}

type GetUsersParams struct {
	UsernamePattern *string
	IDPattern       *string // '123' matches '123456'
	IsBot           *bool

	Offset int
	Limit  int

	OrderBy        string
	OrderDirection string
}

type GetUsersOutput struct {
	Count int
	Users []models.User
}

type UpdateUsersInput struct {
	Name        *string
	Password    *string
	Username    *string
	Description *string
}

type UpdateUsersParams struct {
	Name        *string
	Password    *[]byte
	Username    *string
	Description *string
	AvatarUrl   *string
}
