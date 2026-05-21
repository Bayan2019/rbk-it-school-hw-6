package postgres_test

import (
	"context"
	"testing"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/auth"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_CreateAndGetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewUserRepository(db)

	trueV := true
	createdUser, err := repo.Create(context.Background(), dto.CreateUserInput{
		FirstName: "Alex",
		LastName:  "Some",
		Email:     "alex@example.com",
		IsActive:  &trueV,
	})

	require.NoError(t, err)
	require.NotZero(t, createdUser.ID)

	foundUser, err := repo.GetByEmail(context.Background(), "alex@example.com", true)

	require.NoError(t, err)
	assert.Equal(t, createdUser.FirstName, foundUser.FirstName)
	assert.Equal(t, createdUser.LastName, foundUser.LastName)
	// assert.Equal(t, createdUser.PasswordHash, foundUser.PasswordHash)
	assert.Equal(t, createdUser.Email, foundUser.Email)
}

func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewUserRepository(db)

	tests := []struct {
		name      string
		input     dto.ListUsersFilter
		wantUsers []model.User
		wantErr   error
	}{
		{
			name: "users",
			input: dto.ListUsersFilter{
				Limit:          10,
				Offset:         0,
				Query:          "",
				IncludeDeleted: false,
			},
			wantUsers: []model.User{
				model.User{
					Email:     "user@example.com",
					ID:        2,
					FirstName: "User",
					LastName:  "Role",
					Role:      auth.RolesUser,
				},
				model.User{
					Email:     "admin@example.com",
					ID:        1,
					FirstName: "Admin",
					LastName:  "Role",
					Role:      auth.RolesAdmin,
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users, err := repo.List(context.Background(), tt.input)
			if err != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				assert.Equal(t, len(tt.wantUsers), len(users))
			}
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewUserRepository(db)

	tests := []struct {
		name     string
		userID   int64
		wantUser model.User
		wantErr  error
	}{
		{
			name:   "user",
			userID: int64(2),
			wantUser: model.User{
				Email:     "user@example.com",
				ID:        2,
				FirstName: "User",
				LastName:  "Role",
				Role:      auth.RolesUser,
			},
			wantErr: nil,
		},
		{
			name:   "admin",
			userID: int64(1),
			wantUser: model.User{
				Email:     "admin@example.com",
				ID:        1,
				FirstName: "Admin",
				LastName:  "Role",
				Role:      auth.RolesAdmin,
			},
			wantErr: nil,
		},
		{
			name:     "not found",
			userID:   int64(10),
			wantUser: model.User{},
			wantErr:  model.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByID(context.Background(), tt.userID, false)
			if err != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Email, user.Email)
				assert.Equal(t, tt.wantUser.FirstName, user.FirstName)
				assert.Equal(t, tt.wantUser.LastName, user.LastName)
			}
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewUserRepository(db)

	tests := []struct {
		name     string
		email    string
		wantUser model.User
		wantErr  error
	}{
		{
			name:  "user",
			email: "user@example.com",
			wantUser: model.User{
				Email:     "user@example.com",
				ID:        2,
				FirstName: "User",
				LastName:  "Role",
				Role:      auth.RolesUser,
			},
			wantErr: nil,
		},
		{
			name:  "admin",
			email: "admin@example.com",
			wantUser: model.User{
				Email:     "admin@example.com",
				ID:        1,
				FirstName: "Admin",
				LastName:  "Role",
				Role:      auth.RolesAdmin,
			},
			wantErr: nil,
		},
		{
			name:     "not found",
			email:    "some@example.com",
			wantUser: model.User{},
			wantErr:  model.ErrUserNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := repo.GetByEmail(context.Background(), tt.email, false)
			if err != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				assert.Equal(t, tt.wantUser.ID, user.ID)
				assert.Equal(t, tt.wantUser.Email, user.Email)
				assert.Equal(t, tt.wantUser.FirstName, user.FirstName)
				assert.Equal(t, tt.wantUser.LastName, user.LastName)
			}
		})
	}
}

func TestUserRepository_UpdateAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewUserRepository(db)

	trueV := true
	err := repo.Update(context.Background(), 2, dto.UpdateUserInput{
		FirstName: "Alex",
		LastName:  "Some",
		Email:     "alex@example.com",
		IsActive:  &trueV,
	})

	require.NoError(t, err)

	foundUser, err := repo.GetByID(context.Background(), 2, true)

	require.NoError(t, err)
	assert.Equal(t, "Alex", foundUser.FirstName)
	assert.Equal(t, "Some", foundUser.LastName)
	// assert.Equal(t, createdUser.PasswordHash, foundUser.PasswordHash)
	assert.Equal(t, "alex@example.com", foundUser.Email)
}

func TestUserRepository_DeleteAndGetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewUserRepository(db)

	err := repo.Delete(context.Background(), 2)

	require.NoError(t, err)

	_, err = repo.GetByID(context.Background(), 2, false)

	assert.Equal(t, model.ErrUserNotFound, err)
}
