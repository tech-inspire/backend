package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tech-inspire/backend/auth-service/internal/apperrors"
	"github.com/tech-inspire/backend/auth-service/internal/models"
)

type SessionRepository struct {
	client redis.UniversalClient
}

func NewSessionRepository(client redis.UniversalClient) *SessionRepository {
	return &SessionRepository{client: client}
}

func getUserSessionsHashKey(userID uuid.UUID) string {
	return fmt.Sprintf("user:%s:sessions", userID)
}

func getSessionHashField(sessionID uuid.UUID) string {
	if t := sessionID.Time(); t != 0 {
		// if version of the uuid is 7, then we can use nested time as an id for sorting
		return strconv.FormatInt(int64(t), 10)
	}

	return sessionID.String()
}

func (s SessionRepository) GetUserSession(ctx context.Context, userID, sessionID uuid.UUID) (*models.Session, error) {
	var (
		key   = getUserSessionsHashKey(userID)
		field = getSessionHashField(sessionID)
	)

	data, err := s.client.HGet(ctx, key, field).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, apperrors.ErrSessionNotFound
		}
		return nil, errors.Errorf("redis: get '%s': %w", key, err)
	}

	var schemaSession Session
	err = json.Unmarshal(data, &schemaSession)
	if err != nil {
		return nil, errors.Errorf("unmarshal session from json: %w", err)
	}

	session := models.Session{
		ID:        sessionID,
		UserID:    userID,
		Token:     schemaSession.Token,
		CreatedAt: schemaSession.CreatedAt,
		ExpiresAt: schemaSession.ExpiresAt,
	}
	return &session, nil
}

func (s SessionRepository) DeleteUserSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	key := getUserSessionsHashKey(userID)
	return s.deleteSession(ctx, key, sessionID)
}

func (s SessionRepository) deleteSession(ctx context.Context, key string, sessionID uuid.UUID) error {
	field := getSessionHashField(sessionID)

	err := s.client.HDel(ctx, key, field).Err()
	if err != nil {
		return errors.Errorf("redis: del (key '%s', field '%s'): %w", key, field, err)
	}

	return nil
}

func (s SessionRepository) AddUserSession(ctx context.Context, session models.Session, sessionLimit int) error {
	key := getUserSessionsHashKey(session.UserID)

	schemaSession := Session{
		Token:     session.Token,
		CreatedAt: session.CreatedAt,
		ExpiresAt: session.ExpiresAt,
	}

	return s.addSession(ctx, key, session.ID, schemaSession, sessionLimit)
}

func (s SessionRepository) addSession(ctx context.Context, key string, sessionID uuid.UUID, session Session, sessionLimit int) error {
	sessionCount, err := s.client.HLen(ctx, key).Result()
	if err != nil {
		return errors.Errorf("redis: hlen (key '%s'): %w", key, err)
	}

	if sessionCount >= int64(sessionLimit) {
		keys, err := s.client.HKeys(ctx, key).Result()
		if err != nil {
			return errors.Errorf("redis: hlen (key '%s'): %w", key, err)
		}

		if len(keys) >= sessionLimit {
			sort.Strings(keys) // sessionIDs are sortable

			fieldsToDelete := keys[sessionLimit-1:]
			err = s.client.HDel(ctx, key, fieldsToDelete...).Err()
			if err != nil {
				return errors.Errorf("redis: hdel (key '%s'): %w", key, err)
			}
		}
	}

	data, err := json.Marshal(session)
	if err != nil {
		return errors.Errorf("marshal session into json: %w", err)
	}

	field := getSessionHashField(sessionID)

	err = s.client.HSet(ctx, key, field, data).Err()
	if err != nil {
		return errors.Errorf("redis: hset (key '%s', field '%s'): %w", key, field, err)
	}

	err = s.client.HExpireAt(ctx, key, session.ExpiresAt, field).Err()
	if err != nil {
		return errors.Errorf("redis: hexpireat (key '%s', field '%s'): %w", key, field, err)
	}

	return nil
}
