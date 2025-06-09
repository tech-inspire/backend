package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/tech-inspire/backend/auth-service/internal/apperrors"
	"github.com/tech-inspire/backend/auth-service/internal/clients/mail"
	"github.com/tech-inspire/backend/auth-service/internal/config"
	"github.com/tech-inspire/backend/auth-service/internal/models"
	"github.com/tech-inspire/backend/auth-service/internal/service/dto"
	"github.com/tech-inspire/backend/auth-service/pkg/generator"
	"github.com/tech-inspire/backend/auth-service/pkg/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	logger *logger.Logger

	generator                    Generator
	userRepository               UserRepository
	confirmationCodesRepository  ConfirmationCodesRepository
	resetPasswordCodesRepository ResetPasswordCodesRepository

	sessionRepository SessionRepository
	mailClient        MailClient

	refreshTokenDuration time.Duration
	sessionsLimitPerUser int

	testMode bool
}

func NewAuthService(
	log *logger.Logger,
	cfg *config.Config,

	generator Generator,

	userRepository UserRepository,
	sessionRepository SessionRepository,
	codesRepository ConfirmationCodesRepository,
	resetPasswordCodesRepository ResetPasswordCodesRepository,
	mailClient MailClient,
) *AuthService {
	authService := &AuthService{
		logger: log,

		generator:                    generator,
		userRepository:               userRepository,
		confirmationCodesRepository:  codesRepository,
		resetPasswordCodesRepository: resetPasswordCodesRepository,
		mailClient:                   mailClient,

		sessionRepository: sessionRepository,

		refreshTokenDuration: cfg.JWT.RefreshTokenDuration,
		sessionsLimitPerUser: cfg.Session.MaxAllowedSessionsPerUser,

		testMode: cfg.TestMode,
	}

	log.Info("starting with auth configuration",
		zap.Duration("refresh_token_duration", authService.refreshTokenDuration),
		zap.Int("sessions_limit_per_user", authService.sessionsLimitPerUser),
	)

	return authService
}

func (a AuthService) GetSession(ctx context.Context, userID, sessionID uuid.UUID) (*models.Session, error) {
	session, err := a.sessionRepository.GetUserSession(ctx, userID, sessionID)
	if err != nil {
		return nil, errors.Errorf("create user session: %w", err)
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, apperrors.ErrSessionNotFound
	}

	return session, nil
}

func (a AuthService) DeleteSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	err := a.sessionRepository.DeleteUserSession(ctx, userID, sessionID)
	if err != nil {
		return errors.Errorf("delete user session: %w", err)
	}

	return nil
}

func (a AuthService) Register(ctx context.Context, params dto.RegisterParams) (*dto.RegisterOutput, error) {
	if err := a.checkUsername(ctx, params.Username); err != nil {
		return nil, err
	}
	if err := a.checkMail(ctx, params.Email); err != nil {
		return nil, err
	}

	activeCodesCount, err := a.confirmationCodesRepository.GetActiveCodesCount(ctx, params.Email)
	if err != nil {
		return nil, errors.Errorf("get active codes count: %w", err)
	}

	const maxActiveCodesCount = 5
	if activeCodesCount >= maxActiveCodesCount {
		return nil, errors.Errorf("%w: code was request too many times, try again later", apperrors.ErrForbidden)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Errorf("hash password: %w", err)
	}

	code := "111111"

	if !a.testMode {
		const codeLength = 6
		code = generator.NumberCode(codeLength)

		err = a.mailClient.SendMail(params.Email, mail.ConfirmEmail(code))
		if err != nil {
			return nil, errors.Errorf("send confirmation email: %w", err)
		}
	}

	err = a.confirmationCodesRepository.StoreCode(ctx, models.ConfirmationUserData{
		ConfirmationCode: code,
		Email:            params.Email,
		Username:         params.Username,
		Name:             params.Name,
		PasswordHash:     string(hash),
		ExpiresAt:        time.Now().Add(time.Second * 60 * 5),
	})
	if err != nil {
		return nil, errors.Errorf("store email confirmation code: %w", err)
	}

	return &dto.RegisterOutput{
		ConfirmationRequired: true,
		LoginOutput:          nil,
	}, nil
}

func (a AuthService) checkUsername(ctx context.Context, username string) error {
	_, err := a.userRepository.GetUserByUsername(ctx, username)
	if err == nil {
		return apperrors.ErrUsernameUsed
	}
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		return errors.Errorf("get user by username: %w", err)
	}

	return nil
}

func (a AuthService) checkMail(ctx context.Context, email string) error {
	_, err := a.userRepository.GetUserByEmail(ctx, email)
	if err == nil {
		return apperrors.ErrEmailUsed
	}
	if !errors.Is(err, apperrors.ErrUserNotFound) {
		return errors.Errorf("get user by email: %w", err)
	}

	return nil
}

func (a AuthService) registerUser(ctx context.Context, data models.ConfirmationUserData) (*dto.LoginOutput, error) {
	userID := uuid.Must(uuid.NewV7())

	if err := a.checkUsername(ctx, data.Username); err != nil {
		return nil, err
	}
	if err := a.checkMail(ctx, data.Email); err != nil {
		return nil, err
	}

	err := a.userRepository.CreateUser(ctx, dto.CreateUserParams{
		UserID:       userID,
		Email:        data.Email,
		Name:         data.Name,
		Username:     data.Username,
		PasswordHash: []byte(data.PasswordHash),
		Description:  "",
	})
	if err != nil {
		return nil, errors.Errorf("create user: %w", err)
	}

	sessionID := uuid.Must(uuid.NewV7())
	session, err := a.createSession(ctx, userID, sessionID)
	if err != nil {
		return nil, errors.Errorf("create session: %w", err)
	}

	user, err := a.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, errors.Errorf("get user by id: %w", err)
	}

	return &dto.LoginOutput{
		User:    user,
		Session: session,
	}, nil
}

func (a AuthService) ConfirmRegistrationByCode(ctx context.Context, email, code string) (*dto.LoginOutput, error) {
	data, err := a.confirmationCodesRepository.CheckCode(ctx, email, code)
	if err != nil {
		return nil, errors.Errorf("check code: %w", err)
	}

	return a.registerUser(ctx, *data)
}

func (a AuthService) SendResetPasswordCode(ctx context.Context, email string) error {
	user, err := a.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			// We return nil here so that user couldn't track whether account with such email exists
			return nil
		}

		return fmt.Errorf("get user by mail: %w", err)
	}

	activeCodesCount, err := a.resetPasswordCodesRepository.GetActiveCodesCount(ctx, email)
	if err != nil {
		return errors.Errorf("get active codes count: %w", err)
	}

	const maxActiveCodesCount = 5
	if activeCodesCount >= maxActiveCodesCount {
		return errors.Errorf("%w: code was requested too many times, try again later", apperrors.ErrForbidden)
	}

	const codeLength = 6
	code := "222222" // generator.NumberCode(codeLength)

	err = a.mailClient.SendMail(email, mail.ResetPassword(code))
	if err != nil {
		return errors.Errorf("send confirmation email: %w", err)
	}

	err = a.resetPasswordCodesRepository.StoreCode(ctx, models.ResetPasswordData{
		UserID:    user.ID,
		Code:      code,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(time.Second * 60 * 5),
	})
	if err != nil {
		return errors.Errorf("store email confirmation code: %w", err)
	}

	return nil
}

func (a AuthService) CheckResetPasswordCode(ctx context.Context, email, code string) error {
	_, err := a.resetPasswordCodesRepository.CheckCode(ctx, email, code)
	if err != nil {
		return errors.Errorf("check code: %w", err)
	}
	return nil
}

func (a AuthService) ConfirmResetPasswordByCode(ctx context.Context, email, code, password string) error {
	data, err := a.resetPasswordCodesRepository.CheckCode(ctx, email, code)
	if err != nil {
		return errors.Errorf("check code: %w", err)
	}

	err = a.resetPasswordCodesRepository.DeleteCode(ctx, email, code)
	if err != nil {
		return errors.Errorf("delete code: %w", err)
	}

	_, err = a.userRepository.GetUserByID(ctx, data.UserID)
	if err != nil {
		return fmt.Errorf("get user by id: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.Errorf("hash password: %w", err)
	}

	return a.userRepository.UpdateUserByID(ctx, data.UserID, dto.UpdateUsersParams{
		Password: &hash,
	})
}

func (a AuthService) login(ctx context.Context, user *models.User, passwordHash []byte, password string) (*dto.LoginOutput, error) {
	err := bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
	if err != nil {
		return nil, apperrors.ErrForbidden
	}

	sessionID := uuid.Must(uuid.NewV7())
	session, err := a.createSession(ctx, user.ID, sessionID)
	if err != nil {
		return nil, errors.Errorf("create session: %w", err)
	}

	return &dto.LoginOutput{
		User:    user,
		Session: session,
	}, nil
}

func (a AuthService) LoginByEmail(ctx context.Context, email, password string) (*dto.LoginOutput, error) {
	user, passwordHash, err := a.userRepository.GetUserByEmailWithHash(ctx, email)
	if err != nil {
		return nil, errors.Errorf("get user by email: %w", err)
	}

	return a.login(ctx, user, passwordHash, password)
}

func (a AuthService) LoginByUsername(ctx context.Context, username, password string) (*dto.LoginOutput, error) {
	user, passwordHash, err := a.userRepository.GetUserByUsernameWithHash(ctx, username)
	if err != nil {
		return nil, errors.Errorf("get user by token: %w", err)
	}

	return a.login(ctx, user, passwordHash, password)
}

func (a AuthService) RefreshSession(ctx context.Context, userID, sessionID uuid.UUID, sessionToken string) (*models.User, error) {
	session, err := a.sessionRepository.GetUserSession(ctx, userID, sessionID)
	if err != nil {
		if errors.Is(err, apperrors.ErrSessionNotFound) {
			return nil, apperrors.ErrSessionExpired
		}
		return nil, fmt.Errorf("get session (user '%s', session '%s'): %w", userID, sessionID, err)
	}

	user, err := a.userRepository.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, errors.Errorf("get user by id: %w", err)
	}

	if time.Now().After(session.ExpiresAt) { // in case redis has not deleted it yet
		return nil, apperrors.ErrSessionExpired
	}

	if session.Token != sessionToken {
		return nil, apperrors.ErrForbidden
	}

	return user, nil
}
