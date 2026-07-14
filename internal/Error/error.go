package Error

import "errors"

var (
	ErrFailedConfiguration = errors.New("failed to load configuration")
	ErrUserNotFound        = errors.New("user not found")
	ErrBlogNotFound        = errors.New("blog not found")
)
