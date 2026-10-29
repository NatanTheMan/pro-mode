package service

import (
	"context"
	"strings"
	"unicode/utf8"

	"github.com/NatanTheMan/pro-mode/internal/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, name string) (*domain.Category, error)
}

type CategoryService struct {
	repo CategoryRepository
	adminID,
	devID int64
}

func NewCategoryService(repo CategoryRepository, adminID, devID int64) *CategoryService {
	return &CategoryService{
		repo:    repo,
		adminID: adminID,
		devID:   devID,
	}
}

func (s *CategoryService) Create(ctx context.Context, name string, userID int64) (*domain.Category, error) {
	if userID != s.devID && userID != s.adminID {
		return nil, domain.ErrAccessDenied
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrInvalidName
	}
	if utf8.RuneCountInString(name) < 3 {
		return nil, domain.ErrMinCategoryNameSize
	}
	if utf8.RuneCountInString(name) > 30 {
		return nil, domain.ErrMaxCategoryNameSize
	}

	category, err := s.repo.Create(ctx, name)
	if err != nil {
		return nil, err
	}

	return category, nil
}
