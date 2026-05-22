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

func TestWeatherRepository_CreateHistoryAndHistoryOfUser(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewWeatherRepository(db)

	err := repo.CreateHistory(context.Background(),
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
