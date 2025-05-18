package handlers

import (
	v1 "github.com/tech-inspire/api-contracts/api/gen/go/auth/v1"
	"github.com/tech-inspire/service/auth-service/internal/models"
)

func userPB(u models.User) *v1.User {
	return &v1.User{
		Id: u.ID.String(),
		Username: &v1.Username{
			Value: u.Username,
		},
		Name: &v1.Name{
			Value: u.Name,
		},
		AvatarUrl:   u.AvatarURL,
		Description: u.Description,
	}
}
