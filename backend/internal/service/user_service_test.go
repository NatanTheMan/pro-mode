package service_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/NatanTheMan/pro-mode/internal/domain"
	"github.com/NatanTheMan/pro-mode/internal/service"
	"github.com/google/uuid"
)

type mockRepo struct {
	db []*domain.User
	id uuid.UUID
}

func (r *mockRepo) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	_, err := r.GetByTelegramID(ctx, user.TelegramID)
	if !errors.Is(err, domain.ErrUserNotFound) {
		return nil, err
	}
	if err == nil {
		return nil, domain.ErrUserAlreadyExists
	}

	user.ID = r.id
	r.db = append(r.db, user)

	return user, nil
}

func (r *mockRepo) GetByTelegramID(ctx context.Context, telegramID int64) (*domain.User, error) {
	for _, u := range r.db {
		if u.TelegramID == telegramID {
			return u, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (r *mockRepo) Activate(ctx context.Context, id uuid.UUID, name, username string) (*domain.User, error) {
	for _, u := range r.db {
		if u.ID == id {
			u.Active = true
			u.Name = name
			u.Username = username
			return u, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func TestUserService_Register(t *testing.T) {
	tests := []struct {
		desc string // description of this test case
		repo service.UserRepository
		name,
		username string
		telegramID,
		adminID int64
		want    *domain.User
		wantErr error
	}{
		{
			"success",
			&mockRepo{
				id: uuid.MustParse("91c37bcf-5dd5-4a8a-9009-c227d4e91d1e"),
				db: []*domain.User{},
			},
			"Marco",
			"marco1234lakj342443",
			9128742,
			0,
			&domain.User{
				ID:         uuid.MustParse("91c37bcf-5dd5-4a8a-9009-c227d4e91d1e"),
				Name:       "Marco",
				Username:   "marco1234lakj342443",
				TelegramID: 9128742,
				Active:     true,
				Role:       domain.Member,
			},
			nil,
		},
		{
			"owner with admin rights",
			&mockRepo{
				id: uuid.MustParse("91c37bcf-5dd5-4a8a-9009-c227d4e91d1e"),
				db: []*domain.User{},
			},
			"owner",
			"owner1283749",
			14423889,
			14423889,
			&domain.User{
				ID:         uuid.MustParse("91c37bcf-5dd5-4a8a-9009-c227d4e91d1e"),
				Name:       "owner",
				Username:   "owner1283749",
				TelegramID: 14423889,
				Active:     true,
				Role:       domain.Admin,
			},
			nil,
		},
		{
			"two different users",
			&mockRepo{
				id: uuid.MustParse("91c37bcf-5dd5-4a8a-9009-c227d4e91d1e"),
				db: []*domain.User{
					{
						ID:         uuid.MustParse("9111b01d-975e-443a-bcb7-1878f90a3470"),
						Name:       "malzon",
						Username:   "foo21342",
						TelegramID: 54321,
						Active:     true,
						Role:       domain.Member,
					},
				},
			},
			"john",
			"john1234",
			1234,
			0,
			&domain.User{
				ID:         uuid.MustParse("91c37bcf-5dd5-4a8a-9009-c227d4e91d1e"),
				Name:       "john",
				Username:   "john1234",
				TelegramID: 1234,
				Active:     true,
				Role:       domain.Member,
			},
			nil,
		},
		{
			"active user already exists",
			&mockRepo{
				id: uuid.MustParse("91c37bcf-5dd5-4a8a-9009-c227d4e91d1e"),
				db: []*domain.User{
					{
						ID:         uuid.MustParse("9111b01d-975e-443a-bcb7-1878f90a3470"),
						Name:       "malzon",
						Username:   "foo21342",
						TelegramID: 54321,
						Active:     true,
						Role:       domain.Member,
					},
				},
			},
			"malzon",
			"foo21342",
			54321,
			2,
			nil,
			domain.ErrUserAlreadyExists,
		},
		{
			"activate disabled user",
			&mockRepo{
				id: uuid.MustParse("91c37bcf-5dd5-4a8a-9009-c227d4e91d1e"),
				db: []*domain.User{
					{
						ID:         uuid.MustParse("9111b01d-975e-443a-bcb7-1878f90a3470"),
						Name:       "malzon",
						Username:   "foo21342",
						TelegramID: 54321,
						Active:     false,
						Role:       domain.Member,
					},
				},
			},
			"fisher",
			"chess1234",
			54321,
			2,
			&domain.User{
				ID:         uuid.MustParse("9111b01d-975e-443a-bcb7-1878f90a3470"),
				Name:       "fisher",
				Username:   "chess1234",
				TelegramID: 54321,
				Active:     true,
				Role:       domain.Member,
			},
			nil,
		},
		{
			"empty name",
			&mockRepo{
				id: uuid.MustParse("91c37bcf-5dd5-4a8a-9009-c227d4e91d1e"),
				db: []*domain.User{},
			},
			"",
			"chess1234",
			54321,
			2,
			nil,
			domain.ErrInvalidName,
		},

		{
			"empty username",
			&mockRepo{
				id: uuid.MustParse("91c37bcf-5dd5-4a8a-9009-c227d4e91d1e"),
				db: []*domain.User{},
			},
			"name",
			"",
			54321,
			2,
			&domain.User{
				Name:       "name",
				TelegramID: 54321,
				Role:       domain.Member,
				Active:     true, Username: "",
				ID: uuid.MustParse("91c37bcf-5dd5-4a8a-9009-c227d4e91d1e"),
			},
			nil,
		},

		{
			"name length out of bounds",
			&mockRepo{},
			"ssdflkdsjfsldkfjsdflkdsjfsldkfjdflkdsjfsldkfjsdflkdsjfsldkfjlkjsz",
			"username",
			1234,
			2,
			nil,
			domain.ErrMaxNameSize,
		},

		{
			"username length out of bounds",
			&mockRepo{},
			"name",
			"sldkfowefnonasdofijoierowerfsodf1",
			1234,
			2,
			nil,
			domain.ErrMaxUsernameSize,
		},
		{
			"trim name",
			&mockRepo{},
			"		name     ",
			"username",
			1234,
			2,
			&domain.User{Name: "name", Username: "username", Role: domain.Member, Active: true, TelegramID: 1234},
			nil,
		},
		{
			"spaced name",
			&mockRepo{},
			"   ",
			"username",
			1234,
			2,
			nil,
			domain.ErrInvalidName,
		},

		{
			"fail on active user with bad new details",
			&mockRepo{
				db: []*domain.User{
					{
						ID:         uuid.MustParse("9111b01d-975e-443a-bcb7-1878f90a3470"),
						Name:       "malzon",
						Username:   "foo21342",
						TelegramID: 54321,
						Active:     false,
						Role:       domain.Member,
					},
				},
			},
			" ",
			"chess1234",
			54321,
			2,
			nil,
			domain.ErrInvalidName,
		},
		{
			"name with accents",
			&mockRepo{
				id: uuid.MustParse("9111b01d-975e-443a-bcb7-1878f90a3470"),
			},
			"ááíéáíóúóááíéáíóáíéáíóúóáúóáíéáíóóúóáúóá",
			"username",
			1234,
			2,
			&domain.User{
				ID:         uuid.MustParse("9111b01d-975e-443a-bcb7-1878f90a3470"),
				Name:       "ááíéáíóúóááíéáíóáíéáíóúóáúóáíéáíóóúóáúóá",
				Username:   "username",
				TelegramID: 1234,
				Active:     true,
				Role:       domain.Member,
			},
			nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			s := service.NewUserService(tt.repo, tt.adminID)

			got, gotErr := s.Register(context.Background(), tt.name, tt.username, tt.telegramID)

			if gotErr != nil {
				if !errors.Is(gotErr, tt.wantErr) {
					t.Errorf("Register() failed, got: %v, wantErr: %v", gotErr, tt.wantErr)
				}
				if got != nil {
					t.Errorf("Register() failed but got had: %v	", got)
				}
				return
			}
			if !errors.Is(gotErr, tt.wantErr) {
				t.Fatal("Register() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got: %v, want: %v", got, tt.want)
			}
		})
	}
}
