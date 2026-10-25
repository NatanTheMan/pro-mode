// Package repository
package repository

import (
	"context"
	"errors"

	"github.com/NatanTheMan/pro-mode/internal/domain"
	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByTelegramID(ctx context.Context, telegramID int64) (domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(ctx,
		`SELECT 
				id,
				name,
				COALESCE(username, '') as username,
				role,
				active,
				created_at,
				telegram_id
			FROM users
			WHERE telegram_id = $1;`, telegramID).Scan(
		&u.ID,
		&u.Name,
		&u.Username,
		&u.Role,
		&u.Active,
		&u.CreatedAt,
		&u.TelegramID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return u, nil
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	err := r.db.QueryRow(ctx,
		`INSERT INTO
				users (name, username, role, telegram_id)
			VALUES ($1, NULLIF($2, ''), $3, $4)
			RETURNING id, created_at, active;`,
		user.Name,
		user.Username,
		user.Role,
		user.TelegramID).Scan(&user.ID, &user.CreatedAt, &user.Active)

	var pgErr *pgconn.PgError
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.User{}, domain.ErrUserAlreadyExists
		}
		return domain.User{}, err
	}

	return user, nil
}
