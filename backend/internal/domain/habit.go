package domain

import (
	"time"

	"github.com/google/uuid"
)

type Habit struct {
	ID         uuid.UUID
	Name       string
	Score      int
	Active     bool
	CreatedAt  time.Time
	CategoryID uuid.UUID
}

type InputHabit struct {
	Name       string
	Score      int
	CategoryID uuid.UUID
}
