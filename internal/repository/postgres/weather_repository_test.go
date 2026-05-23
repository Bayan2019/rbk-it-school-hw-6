package postgres_test

import (
	"context"
	"testing"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/repository/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeatherRepository_DoesUserHaveCity_Yes(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewWeatherRepository(db)
	cityRepo := postgres.NewCityRepository(db)

	err := cityRepo.Add2User(context.Background(), 2, dto.AddCityInput{
		City: "Paris",
	})

	require.NoError(t, err)

	exists, err := repo.DoesUserHaveCity(
		context.Background(),
		2,
		"Paris",
	)

	require.NoError(t, err)
	assert.Equal(t, exists, true)
}

func TestWeatherRepository_DoesUserHaveCity_No(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewWeatherRepository(db)

	exists, err := repo.DoesUserHaveCity(
		context.Background(),
		2,
		"Paris",
	)

	require.NoError(t, err)
	assert.Equal(t, exists, false)
}

func TestWeatherRepository_CreateHistoryAndHistoryOfUser(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewWeatherRepository(db)

	err := repo.Create(context.Background(),
		2,
		dto.CityWeatherInput{
			City:        "Paris",
			Temperature: 21.3,
			Description: "sunny",
		})

	require.NoError(t, err)

	expectedHistories := []model.WeatherHistory{
		model.WeatherHistory{
			UserID:      2,
			City:        "Paris",
			Temperature: 21.3,
			Description: "sunny",
		},
	}

	histories, err := repo.WeatherHistoryOfUser(
		context.Background(),
		2,
		dto.WeatherHistoryFilter{
			City:   "",
			Limit:  10,
			Offset: 0,
		},
	)

	expectedHistories[0].RequestedAt = histories[0].RequestedAt

	require.NoError(t, err)
	assert.Equal(t, histories, expectedHistories)

}

func TestWeatherRepository_HistoryOfUser(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewWeatherRepository(db)

	expectedHistories := []model.WeatherHistory{}

	histories, err := repo.WeatherHistoryOfUser(
		context.Background(),
		2,
		dto.WeatherHistoryFilter{
			City:   "",
			Limit:  10,
			Offset: 0,
		},
	)

	require.NoError(t, err)
	assert.Equal(t, expectedHistories, histories)

}
