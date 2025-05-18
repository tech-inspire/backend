package service

import (
	"context"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/tech-inspire/service/auth-service/internal/models"
	"github.com/tech-inspire/service/auth-service/internal/service/dto"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	authService *AuthService

	userRepository UserRepository
}

func NewUserService(authService *AuthService, userRepository UserRepository) *UserService {
	return &UserService{authService: authService, userRepository: userRepository}
}

func (a UserService) GetUserInfoByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := a.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.Errorf("get user by id: %w", err)
	}

	return user, nil
}

func (a UserService) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	_, err := a.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return errors.Errorf("get user '%s': %w", userID, err)
	}

	return a.cleanUpUser(ctx, userID)
}

func (a UserService) UpdateUser(ctx context.Context, userID uuid.UUID, params dto.UpdateUsersInput) error {
	if params.Username != nil {
		if err := a.authService.checkUsername(ctx, *params.Username); err != nil {
			return err
		}
	}

	var passwordHash *[]byte

	if params.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*params.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.Errorf("hash password: %w", err)
		}

		passwordHash = &hash
	}

	return a.userRepository.UpdateUserByID(ctx, userID, dto.UpdateUsersParams{
		Name:        params.Name,
		Password:    passwordHash,
		Username:    params.Username,
		Description: params.Description,
		AvatarUrl:   nil,
	})
}

func (a UserService) GetUserByID(ctx context.Context, userID uuid.UUID) (*dto.GetUserByIDOutput, error) {
	user, err := a.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.Errorf("get user by id: %w", err)
	}

	return &dto.GetUserByIDOutput{
		User: user,
	}, nil
}

func (a UserService) GetCurrentUserByID(ctx context.Context, userID uuid.UUID) (*dto.GetCurrentUser, error) {
	user, err := a.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.Errorf("get user by id: %w", err)
	}

	return &dto.GetCurrentUser{
		GetUserByIDOutput: *user,
	}, nil
}

func (a UserService) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]dto.GetUserByIDOutput, error) {
	users, err := a.userRepository.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, errors.Errorf("get users: %w", err)
	}

	out := make([]dto.GetUserByIDOutput, len(users))
	for i, user := range users {
		out[i].User = user
	}

	return out, nil
}

func (a UserService) GetUsersInfoByID(ctx context.Context, userIDs []uuid.UUID) ([]models.User, error) {
	users, err := a.userRepository.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		return nil, errors.Errorf("get users: %w", err)
	}

	out := make([]models.User, len(users))
	for i, user := range users {
		out[i] = *user
	}

	return out, nil
}

func (a UserService) cleanUpUser(ctx context.Context, userID uuid.UUID) error {
	_, err := a.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return errors.Errorf("get user by id: %w", err)
	}

	if err = a.userRepository.DeleteUserByID(ctx, userID); err != nil {
		return errors.Errorf("delete user: %w", err)
	}

	return nil
}

func (a UserService) GetUsers(ctx context.Context, params dto.GetUsersParams) (*dto.GetUsersOutput, error) {
	users, err := a.userRepository.GetUsers(ctx, params)
	if err != nil {
		return nil, errors.Errorf("get users: %w", err)
	}

	count, err := a.userRepository.GetUsersCount(ctx, params)
	if err != nil {
		return nil, errors.Errorf("count users: %w", err)
	}

	return &dto.GetUsersOutput{
		Count: count,
		Users: users,
	}, nil
}
