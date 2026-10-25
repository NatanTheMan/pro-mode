// Package domain
package domain

import (
	"time"

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
