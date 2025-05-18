package postgres

import (
	"context"
	"fmt"
	"slices"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tech-inspire/service/auth-service/internal/apperrors"
	"github.com/tech-inspire/service/auth-service/internal/models"
	"github.com/tech-inspire/service/auth-service/internal/repository/postgres/sqlc"
	"github.com/tech-inspire/service/auth-service/internal/service/dto"
)

type UserRepository struct {
	repo *sqlc.Queries
	pool *pgxpool.Pool
}

func NewUserRepository(repo *sqlc.Queries, pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{repo: repo, pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, params dto.CreateUserParams) error {
	err := r.repo.CreateUser(ctx, sqlc.CreateUserParams{
		UserID:       params.UserID,
		Email:        params.Email,
		Name:         params.Name,
		Username:     params.Username,
		PasswordHash: params.PasswordHash,
		Description:  params.Description,
	})
	if err != nil {
		return errors.Errorf("sqlc: create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := r.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, errors.Errorf("sqlc: GetUserByUsername: %w", err)
	}

	return userToModel(user.User), nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := r.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, errors.Errorf("sqlc: GetUserByEmail: %w", err)
	}

	return userToModel(user.User), nil
}

func (r *UserRepository) GetUserByUsernameWithHash(ctx context.Context, username string) (admin *models.User, hash []byte, err error) {
	res, err := r.repo.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, apperrors.ErrUserNotFound
		}

		return nil, nil, errors.Errorf("sqlc: GetUserByUsername(%s): %w", username, err)
	}

	return userToModel(res.User), res.User.PasswordHash, nil
}

func (r *UserRepository) GetUserByEmailWithHash(ctx context.Context, email string) (admin *models.User, hash []byte, err error) {
	res, err := r.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, apperrors.ErrUserNotFound
		}

		return nil, nil, errors.Errorf("sqlc: GetUserByEmail(%s): %w", email, err)
	}

	return userToModel(res.User), res.User.PasswordHash, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := r.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, errors.Errorf("sqlc: GetUserByID: %w", err)
	}

	return userToModel(user.User), nil
}

func (r *UserRepository) GetUsersByIDs(ctx context.Context, userIDs []uuid.UUID) ([]*models.User, error) {
	users, err := r.repo.GetUsersByIDs(ctx, userIDs)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, errors.Errorf("sqlc: GetUsersByIDs: %w", err)
	}

	out := make([]*models.User, len(users))
	for i, user := range users {
		out[i] = userToModel(user.User)
	}

	return out, nil
}

func (r *UserRepository) UpdateUserPassword(ctx context.Context, userID uuid.UUID, hash []byte) error {
	err := r.repo.UpdateUserPassword(ctx, hash, userID)
	if err != nil {
		return errors.Errorf("sqlc: UpdateUserPassword: %w", err)
	}

	return nil
}

func (r *UserRepository) DeleteUserByID(ctx context.Context, userID uuid.UUID) error {
	err := r.repo.DeleteUserByID(ctx, userID)
	if err != nil {
		return errors.Errorf("sqlc: DeleteUserByID: %w", err)
	}

	return nil
}

func (*UserRepository) validateOrderField(orderBy string, direction string) (string, error) {
	if !slices.Contains(
		[]string{"name", "created_at"},
		orderBy,
	) {
		return "", errors.Errorf("order field '%s' is not supported", orderBy)
	}

	if !slices.Contains([]string{"desc", "asc"}, direction) {
		return "", errors.Errorf("order direction '%s' is not supported", direction)
	}

	return fmt.Sprintf("%s %s", orderBy, direction), nil
}

func (*UserRepository) applyFilters(builder *sqlbuilder.SelectBuilder, params dto.GetUsersParams) {
	if params.UsernamePattern != nil {
		builder.Where(builder.EQ("username", *params.UsernamePattern+"%"))
	}
	if params.IDPattern != nil {
		builder.Where(builder.Like("user_id::text", *params.IDPattern+"%"))
	}
	if params.IsBot != nil {
		builder.Where(builder.EQ("is_bot", *params.IsBot))
	}
}

func (r *UserRepository) GetUsersCount(ctx context.Context, params dto.GetUsersParams) (count int, err error) {
	builder := sqlbuilder.Select("COUNT(*)").From("users")

	r.applyFilters(builder, params)

	query, args := builder.BuildWithFlavor(sqlbuilder.PostgreSQL)
	err = r.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return -1, errors.Errorf("sqlc: GetUsers: %w", err)
	}

	return count, nil
}

func (r *UserRepository) GetUsers(ctx context.Context, params dto.GetUsersParams) ([]models.User, error) {
	builder := sqlbuilder.Select("*").
		From("users")

	order, err := r.validateOrderField(params.OrderBy, params.OrderDirection)
	if err != nil {
		return nil, err
	}

	r.applyFilters(builder, params)

	builder.OrderBy(order).Offset(params.Offset).Limit(params.Limit)

	query, args := builder.BuildWithFlavor(sqlbuilder.PostgreSQL)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []models.User{}, nil
		}

		return nil, errors.Errorf("pgx: query (args: %v): %w", args, err)
	}

	users, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[sqlc.User])
	if err != nil {
		return nil, errors.Errorf("pgx: collect rows: %w", err)
	}

	out := make([]models.User, len(users))
	for i, user := range users {
		out[i] = *userToModel(user)
	}

	return out, nil
}

func (r *UserRepository) UpdateUserByID(ctx context.Context, userID uuid.UUID, params dto.UpdateUsersParams) error {
	var passwordHash []byte
	if params.Password != nil {
		passwordHash = *params.Password
	}

	return r.repo.UpdateUserByID(ctx, sqlc.UpdateUserByIDParams{
		Name:         params.Name,
		PasswordHash: passwordHash,
		Username:     params.Username,
		Description:  params.Description,
		AvatarUrl:    params.AvatarUrl,
		UserID:       userID,
	})
}

func (r *UserRepository) ClearUserAvatarURL(ctx context.Context, userID uuid.UUID) error {
	return r.repo.ClearUserAvatarURL(ctx, userID)
}
