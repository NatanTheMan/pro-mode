package domain

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidName       = errors.New("invalid name")
	ErrMaxNameSize       = errors.New("name max length is 64")
	ErrMaxUsernameSize   = errors.New("username max length is 32")
)
