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

func TestCategoryRepo_Create(t *testing.T) {
	ctx := context.Background()
	repo := repository.NewCategoryRepository(pool)

	t.Run("inserts and returns category", func(t *testing.T) {
		resetDB(t)

		got, err := repo.Create(ctx, "Saúde")
		if err != nil {
			t.Fatalf("Create() error = %v", err)
		}
		if got.ID == uuid.Nil || got.Name != "Saúde" {
			t.Errorf("Create() = %+v", got)
		}
	})

	t.Run("duplicate name", func(t *testing.T) {
		resetDB(t)

		if _, err := repo.Create(ctx, "Mente"); err != nil {
			t.Fatal(err)
		}
		_, err := repo.Create(ctx, "Mente")
		if !errors.Is(err, domain.ErrCategoryAlreadyExists) {
			t.Errorf("error = %v, want %v", err, domain.ErrCategoryAlreadyExists)
		}
	})

	t.Run("duplicate ignoring case", func(t *testing.T) {
		resetDB(t)

		if _, err := repo.Create(ctx, "Mente"); err != nil {
			t.Fatal(err)
		}
		_, err := repo.Create(ctx, "mente")
		if !errors.Is(err, domain.ErrCategoryAlreadyExists) {
			t.Errorf("error = %v, want %v", err, domain.ErrCategoryAlreadyExists)
		}
	})

	t.Run("concurrent creates, only one wins", func(t *testing.T) {
		resetDB(t)

		const n = 10
		results := make(chan error, n)
		for i := 0; i < n; i++ {
			go func() {
				_, err := repo.Create(ctx, "Foco")
				results <- err
			}()
		}

		var ok, dup int
		for i := 0; i < n; i++ {
			switch err := <-results; {
			case err == nil:
				ok++
			case errors.Is(err, domain.ErrCategoryAlreadyExists):
				dup++
			default:
				t.Errorf("unexpected error: %v", err)
			}
		}
		if ok != 1 || dup != n-1 {
			t.Errorf("ok=%d dup=%d, want 1 and %d", ok, dup, n-1)
		}
	})
}
