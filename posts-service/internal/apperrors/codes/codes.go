package codes

type Code string

const (
	Unauthorized Code = "UNAUTHORIZED"
	Forbidden    Code = "FORBIDDEN"
	PostNotFound Code = "POST_NOT_FOUND"
)
