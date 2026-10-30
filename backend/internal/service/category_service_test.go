package service_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/NatanTheMan/pro-mode/internal/domain"
	"github.com/NatanTheMan/pro-mode/internal/service"
	"github.com/google/uuid"
)

type mockCategoryRepo struct {
	id uuid.UUID
	db []*domain.Category
}

func (r *mockCategoryRepo) Create(ctx context.Context, name string) (*domain.Category, error) {
	for _, c := range r.db {
		if strings.EqualFold(c.Name, name) {
			return nil, domain.ErrCategoryAlreadyExists
		}
	}

	return &domain.Category{ID: r.id, Name: name}, nil
}

func TestCategoryService_Create(t *testing.T) {
	tests := []struct {
		desc string
		repo service.CategoryRepository
		userID,
		adminID int64
		name    string
		want    *domain.Category
		wantErr error
	}{
		{
			"success case",
			&mockCategoryRepo{
				id: uuid.MustParse("21788057-b80e-440f-a44b-56833bbb8aa9"),
			},
			4321,
			4321,
			"		Saude ",
			&domain.Category{ID: uuid.MustParse("21788057-b80e-440f-a44b-56833bbb8aa9"), Name: "Saude"},
			nil,
		},
		{
			"user does not have permission",
			&mockCategoryRepo{},
			123442343,
			1234,
			"Saude",
			nil,
			domain.ErrAccessDenied,
		},
		{
			"invalid name",
			&mockCategoryRepo{},
			1234,
			1234,
			"    			",
			nil,
			domain.ErrInvalidName,
		},
		{
			"name with not enough characters",
			&mockCategoryRepo{},
			1234,
			1234,
			"Ab",
			nil,
			domain.ErrMinCategoryNameSize,
		},
		{
			"name exceeds characters limit",
			&mockCategoryRepo{},
			1234,
			1234,
			"asidufsdna slkdj slkdfjslk sdfa",
			nil,
			domain.ErrMaxCategoryNameSize,
		},
		{
			"category already exists",
			&mockCategoryRepo{
				db: []*domain.Category{
					{ID: uuid.MustParse("feaafa0a-cbcc-4b72-bdcb-6748c5aebaf3"), Name: "Mente"},
				},
			},
			1234,
			1234,
			"mente",
			nil,
			domain.ErrCategoryAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			s := service.NewCategoryService(tt.repo, tt.adminID)
			got, gotErr := s.Create(context.Background(), tt.name, tt.userID)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Fatalf("Create() error = %v, want %v", gotErr, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
