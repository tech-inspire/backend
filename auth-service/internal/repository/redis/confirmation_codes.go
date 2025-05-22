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

type ConfirmationCodesRepository struct {
	client redis.UniversalClient
}

func NewCodesRepository(client redis.UniversalClient) *ConfirmationCodesRepository {
	return &ConfirmationCodesRepository{
		client: client,
	}
}

func (*ConfirmationCodesRepository) getKey(email string) string {
	return fmt.Sprintf("codes:confirmation:%s", email)
}

func (repo *ConfirmationCodesRepository) StoreCode(ctx context.Context, data models.ConfirmationUserData) error {
	confirmationData := UserConfirmationData(data)

	bytes, err := json.Marshal(confirmationData)
	if err != nil {
		return fmt.Errorf("marshal user data: %w", err)
	}

	key := repo.getKey(data.Email)
	field := data.ConfirmationCode

	err = repo.client.HSet(ctx, key, field, bytes).Err()
	if err != nil {
		return fmt.Errorf("redis: set field for key %s: %w", key, err)
	}

	err = repo.client.HExpireAt(ctx, key, confirmationData.ExpiresAt, field).Err()
	if err != nil {
		return fmt.Errorf("redis: set expiration for key %s field %s: %w", key, field, err)
	}

	return nil
}

func (repo *ConfirmationCodesRepository) GetActiveCodesCount(ctx context.Context, email string) (int64, error) {
	key := repo.getKey(email)

	count, err := repo.client.HLen(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("redis: hlen: key %s: %w", email, err)
	}

	return count, nil
}

func (repo *ConfirmationCodesRepository) ClearAllCodes(ctx context.Context, email string) error {
	key := repo.getKey(email)

	err := repo.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("redis: del: key %s: %w", email, err)
	}

	return nil
}

func (repo *ConfirmationCodesRepository) CheckCode(ctx context.Context, email, confirmationCode string) (*models.ConfirmationUserData, error) {
	key := repo.getKey(email)
	field := confirmationCode

	data, err := repo.client.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, apperrors.ErrConfirmationCodeNotFound
		}

		return nil, fmt.Errorf("redis: get field %s for key %s: %w", confirmationCode, key, err)
	}

	var confirmationData UserConfirmationData
	err = json.Unmarshal([]byte(data), &confirmationData)
	if err != nil {
		return nil, fmt.Errorf("unmarshal data for field %s in key %s: %w", confirmationCode, key, err)
	}

	if confirmationData.ConfirmationCode == confirmationCode && time.Now().Before(confirmationData.ExpiresAt) {
		err = repo.client.HDel(ctx, key, confirmationCode).Err()
		if err != nil {
			return nil, fmt.Errorf("redis: delete field %s for key %s: %w", confirmationCode, key, err)
		}

		user := models.ConfirmationUserData(confirmationData)

		return &user, nil
	}

	// FIXME: the hell we return nil?

	return nil, nil
}
