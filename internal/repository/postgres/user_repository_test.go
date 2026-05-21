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
