package service

import (
	"context"
	"strings"

	"github.com/NatanTheMan/pro-mode/internal/domain"
	"github.com/google/uuid"
)

type HabitRepository interface {
	Create(ctx context.Context, habit *domain.InputHabit) (*domain.Habit, error)
}

type HabitService struct {
	repo    HabitRepository
	adminID int64
}

func NewHabitService(repo HabitRepository, adminID int64) *HabitService {
	return &HabitService{repo: repo, adminID: adminID}
}

func (s *HabitService) Create(ctx context.Context, habit *domain.InputHabit, userID int64) (*domain.Habit, error) {
	habit.Name = strings.TrimSpace(habit.Name)

	if userID != s.adminID {
		return nil, domain.ErrAccessDenied
	}
	if habit.Score <= 0 || habit.Score > 100 {
		return nil, domain.ErrInvalidHabitScore
	}
	if habit.CategoryID == uuid.Nil {
		return nil, domain.ErrInvalidCategory
	}
	if habit.Name == "" {
		return nil, domain.ErrInvalidName
	}

	h, err := s.repo.Create(ctx, habit)
	if err != nil {
		return nil, err
	}

	return h, nil
}
