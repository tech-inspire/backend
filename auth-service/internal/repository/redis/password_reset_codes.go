package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tech-inspire/backend/auth-service/internal/apperrors"
	"github.com/tech-inspire/backend/auth-service/internal/models"
)

type ResetPasswordCodesRepository struct {
	client redis.UniversalClient
}

func NewResetCodesRepository(client redis.UniversalClient) *ResetPasswordCodesRepository {
	return &ResetPasswordCodesRepository{
		client: client,
	}
}

func (*ResetPasswordCodesRepository) getKey(email string) string {
	return fmt.Sprintf("codes:password_reset:%s", email)
}

func (repo *ResetPasswordCodesRepository) StoreCode(ctx context.Context, data models.ResetPasswordData) error {
	redisUser := ResetPasswordData{
		UserID:    data.UserID,
		Code:      data.Code,
		ExpiresAt: data.ExpiresAt,
	}

	bytes, err := json.Marshal(redisUser)
	if err != nil {
		return fmt.Errorf("marshal user data: %w", err)
	}

	key := repo.getKey(data.Email)
	field := data.Code

	err = repo.client.HSet(ctx, key, field, bytes).Err()
	if err != nil {
		return fmt.Errorf("redis: set field for key %s: %w", key, err)
	}

	err = repo.client.HExpireAt(ctx, key, redisUser.ExpiresAt, field).Err()
	if err != nil {
		return fmt.Errorf("redis: set expiration for key %s field %s: %w", key, field, err)
	}

	return nil
}

func (repo *ResetPasswordCodesRepository) GetActiveCodesCount(ctx context.Context, email string) (int64, error) {
	key := repo.getKey(email)

	count, err := repo.client.HLen(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis: hlen: key %s: %w", email, err)
	}

	return count, nil
}

func (repo *ResetPasswordCodesRepository) ClearAllCodes(ctx context.Context, email string) error {
	key := repo.getKey(email)

	err := repo.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis: del: key %s: %w", email, err)
	}

	return nil
}

func (repo *ResetPasswordCodesRepository) DeleteCode(ctx context.Context, email, code string) error {
	key := repo.getKey(email)
	field := code

	err := repo.client.HDel(ctx, key, field).Err()
	if err != nil {
		return fmt.Errorf("redis: hdel: key %s: %w", email, err)
	}

	return nil
}

func (repo *ResetPasswordCodesRepository) CheckCode(ctx context.Context, email, confirmationCode string) (*models.ResetPasswordData, error) {
	key := repo.getKey(email)
	field := confirmationCode

	data, err := repo.client.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, apperrors.ErrResetPasswordCodeNotFound
		}

		return nil, fmt.Errorf("redis: get field %s for key %s: %w", confirmationCode, key, err)
	}

	var rsData ResetPasswordData
	err = json.Unmarshal([]byte(data), &rsData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal data for field %s in key %s: %w", confirmationCode, key, err)
	}

	if rsData.Code != confirmationCode || time.Now().After(rsData.ExpiresAt) {
		return nil, apperrors.ErrResetPasswordCodeNotFound
	}

	return &models.ResetPasswordData{
		UserID:    rsData.UserID,
		Code:      rsData.Code,
		Email:     email,
		ExpiresAt: rsData.ExpiresAt,
	}, nil
}
