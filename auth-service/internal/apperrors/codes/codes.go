package codes

type Code string

const (
	Unauthorized Code = "UNAUTHORIZED"

	Forbidden       Code = "FORBIDDEN"
	SessionExpired  Code = "SESSION_EXPIRED"
	SessionNotFound Code = "SESSION_NOT_FOUND"

	UserNotFound Code = "USER_NOT_FOUND"
	EmailUsed    Code = "EMAIL_USED"
	UsernameUsed Code = "USERNAME_USED"

	ConfirmationCodeNotFound  = "CONFIRMATION_CODE_NOT_FOUND"
	ResetPasswordCodeNotFound = "RESET_CODE_NOT_FOUND"
)
