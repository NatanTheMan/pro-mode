// Package domain
package domain

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Role string

const (
	Admin  Role = "ADMIN"
	Member Role = "MEMBER"
)

type User struct {
	ID         uuid.UUID
	Name       string
	Username   string
	Role       Role
	Active     bool
	CreatedAt  time.Time
	TelegramID int64
}

func NewUser(name, username string, telegramID int64) (*User, error) {
	user := &User{
		Name:       strings.TrimSpace(name),
		Username:   strings.TrimSpace(username),
		Role:       Member,
		TelegramID: telegramID,
		Active:     true,
	}

	if user.Name == "" || user.Name == " " {
		return nil, ErrInvalidName
	}
	if utf8.RuneCountInString(user.Name) > 64 {
		return nil, ErrMaxNameSize
	}
	if utf8.RuneCountInString(user.Username) > 32 {
		return nil, ErrMaxUsernameSize
	}

	return user, nil
}
