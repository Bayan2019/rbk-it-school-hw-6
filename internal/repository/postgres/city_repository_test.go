package postgres_test

import (
	"context"
	"testing"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCityRepository_CreateAndGetByName(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewCityRepository(db)

	err := repo.Create(context.Background(), dto.CreateCityInput{
		City: "London",
		Lat:  51.5074456,
		Lon:  -0.1277653,
	})

	require.NoError(t, err)

	foundCity, err := repo.GetByName(context.Background(), "London")

	require.NoError(t, err)
	assert.Equal(t, "London", foundCity.City)
	assert.Equal(t, 51.5074456, foundCity.Lat)
	assert.Equal(t, -0.1277653, foundCity.Lon)
}
