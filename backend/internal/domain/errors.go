package domain

import "errors"

var (
	ErrAccessDenied      = errors.New("access denied")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidName       = errors.New("invalid name")
	ErrMaxNameSize       = errors.New("name max length is 64")
	ErrMaxUsernameSize   = errors.New("username max length is 32")

	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrMinCategoryNameSize   = errors.New("category name should have at least 3 characters")
	ErrMaxCategoryNameSize   = errors.New("category name must have a max of 30 characters")
)
