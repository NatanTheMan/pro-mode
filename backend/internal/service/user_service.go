// Package service
package service

import (
	"context"
	"errors"

	"github.com/NatanTheMan/pro-mode/internal/domain"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error)
	Activate(ctx context.Context, id uuid.UUID, name, username string) (*domain.User, error)
}

type UserService struct {
	repo UserRepository
	adminID,
	devID int64
}

func NewUserService(repo UserRepository, adminID, devID int64) *UserService {
	return &UserService{repo: repo, adminID: adminID, devID: devID}
}

func (s *UserService) Register(ctx context.Context, name, username string, telegramID int64) (*domain.User, error) {
	response, err := s.repo.GetByTelegramID(ctx, telegramID)
	if err == nil {
		if response.Active {
			return nil, domain.ErrUserAlreadyExists
		}
		u, err := domain.NewUser(name, username, telegramID)
		if err != nil {
			return nil, err
		}
		return s.repo.Activate(ctx, response.ID, u.Name, u.Username)
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}

	user, err := domain.NewUser(name, username, telegramID)
	if err != nil {
		return nil, err
	}
	if telegramID == s.adminID || telegramID == s.devID {
		user.Role = domain.Admin
	}
	return s.repo.Create(ctx, user)
}
