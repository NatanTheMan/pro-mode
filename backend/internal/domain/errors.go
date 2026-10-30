package domain

import "errors"

var (
	ErrInvalidName       = errors.New("invalid name")
	ErrMaxNameSize       = errors.New("name max length is 64")
	ErrAccessDenied      = errors.New("access denied")
	ErrUserNotFound      = errors.New("user not found")
	ErrMaxUsernameSize   = errors.New("username max length is 32")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrCategoryNotFound      = errors.New("category not found")
	ErrMinCategoryNameSize   = errors.New("category name should have at least 3 characters")
	ErrMaxCategoryNameSize   = errors.New("category name must have a max of 30 characters")
	ErrCategoryAlreadyExists = errors.New("category already exists")

	ErrInvalidCategory   = errors.New("invalid category uuid")
	ErrInvalidHabitScore = errors.New("invalid habit score")
)
