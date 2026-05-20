package postgres_test

import (
	"context"
	"testing"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSQLUserRepository_CreateAndGetByID(t *testing.T) {
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
