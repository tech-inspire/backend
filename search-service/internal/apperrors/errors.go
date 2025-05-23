package apperrors

import (
	"github.com/tech-inspire/backend/search-service/internal/apperrors/codes"
)

var (
	ErrUnauthorized = newError(codes.Unauthorized, "unauthorized")
	ErrForbidden    = newError(codes.Forbidden, "forbidden")

	ErrSessionExpired = newError(codes.SessionExpired, "session expired")

	ErrUserNotFound    = newError(codes.UserNotFound, "user not found")
	ErrSessionNotFound = newError(codes.SessionNotFound, "session not found")

	ErrEmailUsed    = newError(codes.EmailUsed, "email already used")
	ErrUsernameUsed = newError(codes.UsernameUsed, "username already used")

	ErrConfirmationCodeNotFound  = newError(codes.ConfirmationCodeNotFound, "confirmation code not found")
	ErrResetPasswordCodeNotFound = newError(codes.ResetPasswordCodeNotFound, "reset password code not found")
)
