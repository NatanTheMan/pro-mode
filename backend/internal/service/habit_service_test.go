package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/NatanTheMan/pro-mode/internal/domain"
	"github.com/google/uuid"
)

type mockHabitRepo struct {
	id         uuid.UUID
	categoryID uuid.UUID
	err        error // forced error, e.g. a DB failure
	calls      int
	gotInput   domain.InputHabit
}

func (r *mockHabitRepo) Create(ctx context.Context, in *domain.InputHabit) (*domain.Habit, error) {
	r.calls++
	r.gotInput = *in
	if r.err != nil {
		return nil, r.err
	}
	if in.CategoryID != r.categoryID {
		return nil, domain.ErrCategoryNotFound
	}
	return &domain.Habit{
		ID:         r.id,
		Name:       in.Name,
		Score:      in.Score,
		CategoryID: r.categoryID,
		Active:     true,
	}, nil
}

var categoryID uuid.UUID = uuid.MustParse("2d11c929-dd31-4716-a77a-cb1d2d75f36a")

func TestHabitService_Create(t *testing.T) {
	tests := []struct {
		name string
		repo *mockHabitRepo
		userID,
		adminID int64
		habit     *domain.InputHabit
		want      *domain.Habit
		wantErr   error
		wantCalls int
	}{
		{
			"score zero",
			&mockHabitRepo{},
			1234, 1234,
			&domain.InputHabit{Name: "valid", Score: 0, CategoryID: categoryID},
			nil,
			domain.ErrInvalidHabitScore,
			0,
		},
		{
			"negative points",
			&mockHabitRepo{},
			1234, 1234,
			&domain.InputHabit{
				Name:       "Workout",
				Score:      -1,
				CategoryID: categoryID,
			},
			nil,
			domain.ErrInvalidHabitScore,
			0,
		},
		{
			"empty category id",
			&mockHabitRepo{},
			1234, 1234,
			&domain.InputHabit{Score: 1},
			nil,
			domain.ErrInvalidCategory,
			0,
		},
		{
			"reject not admin",
			&mockHabitRepo{},
			1234, 12345,
			&domain.InputHabit{
				Score:      1,
				CategoryID: categoryID,
			},
			nil,
			domain.ErrAccessDenied,
			0,
		},
		{
			"invalid name",
			&mockHabitRepo{},
			1234, 1234,
			&domain.InputHabit{
				Name:       "",
				Score:      1,
				CategoryID: categoryID,
			},
			nil,
			domain.ErrInvalidName,
			0,
		},
		{
			"blank name",
			&mockHabitRepo{},
			1234, 1234,
			&domain.InputHabit{
				Name:       " 		 ",
				Score:      1,
				CategoryID: categoryID,
			},
			nil,
			domain.ErrInvalidName,
			0,
		},
		{
			"invalid category, category not exists",
			&mockHabitRepo{},
			1234, 1234,
			&domain.InputHabit{
				Name:       "name",
				Score:      2,
				CategoryID: categoryID,
			},
			nil,
			domain.ErrCategoryNotFound,
			1,
		},
		{
			"success case",
			&mockHabitRepo{
				id:         uuid.MustParse("821dc6b4-3b19-4f5a-af15-825cb0039387"),
				categoryID: categoryID,
			},
			1234,
			1234,
			&domain.InputHabit{
				Name:       "name",
				Score:      3,
				CategoryID: categoryID,
			},
			&domain.Habit{
				ID:         uuid.MustParse("821dc6b4-3b19-4f5a-af15-825cb0039387"),
				Name:       "name",
				Score:      3,
				CategoryID: categoryID,
				Active:     true,
			},
			nil,
			1,
		},
		{
			"score above limit",
			&mockHabitRepo{
				id:         uuid.MustParse("821dc6b4-3b19-4f5a-af15-825cb0039387"),
				categoryID: categoryID,
			},
			1234, 1234,
			&domain.InputHabit{Score: 101},
			nil,
			domain.ErrInvalidHabitScore,
			0,
		},
		{
			"score at limit",
			&mockHabitRepo{
				id:         uuid.MustParse("821dc6b4-3b19-4f5a-af15-825cb0039387"),
				categoryID: categoryID,
			},
			1234, 1234,
			&domain.InputHabit{Score: 100, Name: "habit", CategoryID: categoryID},
			&domain.Habit{
				ID:         uuid.MustParse("821dc6b4-3b19-4f5a-af15-825cb0039387"),
				Name:       "habit",
				Score:      100,
				CategoryID: categoryID,
				Active:     true,
			},
			nil,
			1,
		},
		{
			"name with spaces",
			&mockHabitRepo{
				id:         uuid.MustParse("821dc6b4-3b19-4f5a-af15-825cb0039387"),
				categoryID: categoryID,
			},
			1234, 1234,
			&domain.InputHabit{Score: 100, Name: "  		habit 		", CategoryID: categoryID},
			&domain.Habit{
				ID:         uuid.MustParse("821dc6b4-3b19-4f5a-af15-825cb0039387"),
				Name:       "habit",
				Score:      100,
				CategoryID: categoryID,
				Active:     true,
			},
			nil,
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewHabitService(tt.repo, tt.adminID)
			got, gotErr := s.Create(context.Background(), tt.habit, tt.userID)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Fatalf("Create() error = %v, want %v", gotErr, tt.wantErr)
			}
			if tt.repo.calls != tt.wantCalls {
				t.Errorf("repo calls = %d, want %d", tt.repo.calls, tt.wantCalls)
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
