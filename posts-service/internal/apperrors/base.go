package apperrors

import (
	"errors"

	"github.com/tech-inspire/backend/posts-service/internal/apperrors/codes"
)

var _ error = &Error{}

type Error struct {
	Code codes.Code
	Err  error
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func newError(code codes.Code, err string) *Error {
	return &Error{
		Code: code,
		Err:  errors.New(err),
	}
}
