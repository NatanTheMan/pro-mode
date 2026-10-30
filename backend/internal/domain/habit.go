package domain

import (
	"time"

	"github.com/google/uuid"
)

type Habit struct {
	ID         uuid.UUID
	Name       string
	Points     int
	Active     bool
	CreatedAt  time.Time
	CategoryID uuid.UUID
}
