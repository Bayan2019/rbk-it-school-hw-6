package service_test

import (
	"context"
	"testing"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/client"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/config"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

/// WeatherRepository ///
/// WeatherRepository ///
/// WeatherRepository ///
/// WeatherRepository ///
/// WeatherRepository ///

type MockWeatherRepository struct {
	mock.Mock
}

func (m *MockWeatherRepository) DoesUserHaveCity(
	ctx context.Context,
	userID int64,
	city string,
) (bool, error) {
	args := m.Called(ctx, userID, city)

	exists, _ := args.Get(0).(bool)
	return exists, args.Error(1)
}

func (m *MockWeatherRepository) Create(
	ctx context.Context,
	userID int64,
	input dto.CityWeatherInput,
) error {
	args := m.Called(ctx, userID, input)

	return args.Error(0)
}

func (m *MockWeatherRepository) WeatherHistoryOfUser(
	ctx context.Context,
	userID int64,
	filter dto.WeatherHistoryFilter,
) ([]model.WeatherHistory, error) {
	args := m.Called(ctx, userID, filter)

	cities, _ := args.Get(0).([]model.WeatherHistory)
	return cities, args.Error(1)
}

/// WeatherServiceTests ///
/// WeatherServiceTests ///
/// WeatherServiceTests ///
/// WeatherServiceTests ///
/// WeatherServiceTests ///

func TestWeatherService_Create_Success(t *testing.T) {
	err := config.MustLoad("../../.env")
	require.NoError(t, err)

	client := client.NewWeatherClient()

	repo := new(MockWeatherRepository)
	weatherService := service.NewWeatherService(repo, client)

	cityInput := model.City{
		CityID: 1,
		City:   "Paris",
		Lat:    48.8534951,
		Lon:    2.3483915,
	}

	repo.On("DoesUserHaveCity", mock.Anything, int64(2), cityInput.City).
		Return(true, nil).
		Once()

	repo.On("Create", mock.Anything, int64(2), mock.MatchedBy(func(input dto.CityWeatherInput) bool {
		return input.City == "Paris"
	})).
		Return(nil).
		Once()

	_, err = weatherService.Create(context.Background(), 2, cityInput)

	require.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestWeatherService_WeatherHistoryOfUser(t *testing.T) {
	err := config.MustLoad("../../.env")
	require.NoError(t, err)

	client := client.NewWeatherClient()

	repo := new(MockWeatherRepository)
	weatherService := service.NewWeatherService(repo, client)

	filter := dto.WeatherHistoryFilter{
		City:   "",
		Limit:  10,
		Offset: 0,
	}

	repo.On("WeatherHistoryOfUser", mock.Anything, int64(2), filter).
		Return([]model.WeatherHistory{}, nil).
		Once()

	actualResult, err := weatherService.WeatherHistoryOfUser(context.Background(), 2, filter)

	require.NoError(t, err)
	assert.Equal(t, []model.WeatherHistory{}, actualResult)

	repo.AssertExpectations(t)
}
