package repository

import (
	"context"
	"errors"

	"github.com/NatanTheMan/pro-mode/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, name string) (*domain.Category, error) {
	var c domain.Category
	err := r.db.QueryRow(ctx,
		"INSERT INTO categories(name) VALUES ($1) RETURNING id, name;",
		name).Scan(&c.ID, &c.Name)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrCategoryAlreadyExists
		}
		return nil, err
	}

	return &c, nil
}
