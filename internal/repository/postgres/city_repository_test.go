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

func TestCityRepository_Create_AlreadyExist(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewCityRepository(db)

	err := repo.Create(context.Background(), dto.CreateCityInput{
		City: "Paris",
		Lat:  48.8534951,
		Lon:  2.3483915,
	})

	assert.Equal(t, model.ErrCityNameAlreadyTaken, err)
}

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

func TestCityRepository_Add2UserAndListCitiesOfUser(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewCityRepository(db)

	err := repo.Add2User(context.Background(), 2, dto.AddCityInput{
		City: "Paris",
	})

	require.NoError(t, err)

	cities, err := repo.ListCitiesOfUser(
		context.Background(),
		2,
		dto.ListCitiesFilter{
			Offset:         0,
			IncludeDeleted: false,
		},
	)

	require.NoError(t, err)
	city := cities[0]
	assert.Equal(t, "Paris", city.City)
}

func TestCityRepository_Add2User_AlreadyAdded(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewCityRepository(db)

	err := repo.Add2User(context.Background(), 2, dto.AddCityInput{
		City: "Paris",
	})

	require.NoError(t, err)

	err = repo.Add2User(
		context.Background(),
		2,
		dto.AddCityInput{
			City: "Paris",
		},
	)

	assert.Equal(t, model.ErrCityAlreadyAdded2User, err)
}

func TestCityRepository_GetByName(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewCityRepository(db)

	tests := []struct {
		name     string
		city     string
		wantCity model.City
		wantErr  error
	}{
		{
			name: "Paris",
			city: "Paris",
			wantCity: model.City{
				CityID: 1,
				City:   "Paris",
			},
			wantErr: nil,
		},
		{
			name: "Berlin",
			city: "Berlin",
			wantCity: model.City{
				CityID: 2,
				City:   "Berlin",
			},
			wantErr: nil,
		},
		{
			name:     "not found",
			city:     "London",
			wantCity: model.City{},
			wantErr:  model.ErrCityNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			city, err := repo.GetByName(context.Background(), tt.city)
			if err != tt.wantErr {
				t.Errorf("GetByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				assert.Equal(t, tt.wantCity.CityID, city.CityID)
				assert.Equal(t, tt.wantCity.City, city.City)
			}
		})
	}
}

func TestCityRepository_DeleteFromUser(t *testing.T) {
	db := setupTestDB(t)
	repo := postgres.NewCityRepository(db)

	err := repo.DeleteFromUser(context.Background(), 2, 1)

	assert.Equal(t, model.ErrCityNotFound, err)
}
