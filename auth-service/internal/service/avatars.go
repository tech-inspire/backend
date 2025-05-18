package service

import (
	"context"
	"fmt"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/tech-inspire/service/auth-service/internal/service/dto"
)

type AvatarService struct {
	userRepository UserRepository
	storage        AvatarStorage
}

func NewAvatarService(userRepository UserRepository, storage AvatarStorage) *AvatarService {
	return &AvatarService{userRepository: userRepository, storage: storage}
}

func (g AvatarService) UploadUserAvatar(ctx context.Context, params dto.UploadUserAvatar) error {
	_, err := g.userRepository.GetUserByID(ctx, params.UserID)
	if err != nil {
		return fmt.Errorf("get user: %w", err)
	}

	path, err := g.storage.UploadUserAvatar(ctx, params)
	if err != nil {
		return errors.Errorf("storage: get upload url: %w", err)
	}

	err = g.userRepository.UpdateUserByID(ctx, params.UserID, dto.UpdateUsersParams{
		AvatarUrl: &path,
	})
	if err != nil {
		return errors.Errorf("update user by id: %w", err)
	}

	return nil
}

func (g AvatarService) DeleteProfileAvatar(ctx context.Context, userID uuid.UUID) error {
	err := g.storage.DeleteUserAvatar(ctx, userID)
	if err != nil {
		return errors.Errorf("storage: delete avatar: %w", err)
	}

	err = g.userRepository.ClearUserAvatarURL(ctx, userID)
	if err != nil {
		return errors.Errorf("clear avatar url: %w", err)
	}

	return nil
}
