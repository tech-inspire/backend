package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenUse string

const (
	AccessToken  TokenUse = "access"
	RefreshToken TokenUse = "refresh"
)

type UserAccessTokenClaims struct {
	jwt.RegisteredClaims
	TokenUse TokenUse `json:"token_use"`

	IsAdmin bool `json:"is_admin,omitempty"`

	SessionID uuid.UUID `json:"session_id"`
}

const (
	Issuer   = "inspire-auth"
	Audience = "inspire-web"
)

type UserRefreshTokenClaims struct {
	jwt.RegisteredClaims
	TokenUse TokenUse `json:"token_use"`

	SessionID    uuid.UUID `json:"session_id"`
	SessionToken string    `json:"session_token"`
}
