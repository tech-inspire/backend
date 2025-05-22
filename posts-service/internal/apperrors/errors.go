package apperrors

import (
	"github.com/tech-inspire/backend/posts-service/internal/apperrors/codes"
)

var (
	ErrUnauthorized = newError(codes.Unauthorized, "unauthorized")
	ErrForbidden    = newError(codes.Forbidden, "forbidden")

	ErrPostNotFound = newError(codes.PostNotFound, "user not found")
)
