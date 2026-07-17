package Error

import "errors"

var (
	ErrFailedConfiguration = errors.New("failed to load configuration")
	ErrUserNotFound        = errors.New("user not found")
	ErrBlogNotFound        = errors.New("blog not found")
	ErrAllFieldRequired    = errors.New("all fields are required")
	ErrIncorrectPassword   = errors.New("incorrect password")
	ErrPasswordTooShort    = errors.New("password must be at least 8 characters")
	ErrInvalidEmail        = errors.New("invalid email format")
	ErrEmailAlreadyExists  = errors.New("email already registered")
	ErrMissingIdentifiers  = errors.New("blog and user id are required")
)
