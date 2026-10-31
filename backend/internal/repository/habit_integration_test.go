//go:build integration

package repository_test

import (
	"context"
	"errors"
	"testing"

	"github.com/NatanTheMan/pro-mode/internal/domain"
	"github.com/NatanTheMan/pro-mode/internal/repository"
	"github.com/google/uuid"
)

func TestHabitRepo_Create(t *testing.T) {
	ctx := context.Background()
	catRepo := repository.NewCategoryRepository(pool)
	sut := repository.NewHabitRepository(pool)

	t.Run("inserts and return habit", func(t *testing.T) {
		resetDB(t)
		category, _ := catRepo.Create(ctx, "Health")

		got, err := sut.Create(ctx, &domain.InputHabit{
			Name:       "Train",
			Score:      2,
			CategoryID: category.ID,
		})
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if got.ID == uuid.Nil ||
			got.Name != "Train" ||
			got.Score != 2 ||
			!got.Active ||
			got.CategoryID != category.ID {
			t.Fatalf("Create() = %+v", got)
		}
	})

	t.Run("duplicated habit", func(t *testing.T) {
		resetDB(t)
		category, _ := catRepo.Create(ctx, "Health")

		habit := &domain.InputHabit{
			Name:       "Train",
			Score:      2,
			CategoryID: category.ID,
		}

		sut.Create(ctx, habit)
		got, err := sut.Create(ctx, habit)

		if !errors.Is(err, domain.ErrHabitAlreadyExists) {
			t.Fatalf("unexpected error = %v", err)
		}
		if got != nil {
			t.Fatalf("unexpected success, = %+v", got)
		}
	})

	t.Run("not existing category", func(t *testing.T) {
		resetDB(t)

		habit := &domain.InputHabit{
			Name:  "Train",
			Score: 2,
		}

		got, err := sut.Create(ctx, habit)

		if err == nil || got != nil {
			t.Fatal("unexpected success")
		}
		if !errors.Is(err, domain.ErrCategoryNotFound) {
			t.Fatalf("unexpected error, want %v got %v", domain.ErrCategoryNotFound, err)
		}
	})
}
