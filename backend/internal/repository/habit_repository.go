package repository

import (
	"context"
	"errors"
	"time"

	"github.com/NatanTheMan/pro-mode/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HabitRepository struct {
	db *pgxpool.Pool
}

func NewHabitRepository(db *pgxpool.Pool) *HabitRepository {
	return &HabitRepository{db: db}
}

func (r *HabitRepository) Create(ctx context.Context, habit *domain.InputHabit) (*domain.Habit, error) {
	var h domain.Habit
	err := r.db.QueryRow(ctx,
		"INSERT INTO habits(name, score, category_id, created_at) VALUES ($1, $2, $3, $4) RETURNING id, name, score, category_id, active, created_at;",
		habit.Name, habit.Score, habit.CategoryID, time.Now()).Scan(
		&h.ID,
		&h.Name,
		&h.Score,
		&h.CategoryID,
		&h.Active,
		&h.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrHabitAlreadyExists
		}
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, domain.ErrCategoryNotFound
		}
	}

	return &h, nil
}
