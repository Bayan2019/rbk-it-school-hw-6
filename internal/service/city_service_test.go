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

/// CityRepository ///
/// CityRepository ///
/// CityRepository ///
/// CityRepository ///
/// CityRepository ///

type MockCityRepository struct {
	mock.Mock
}

func (m *MockCityRepository) Create(
	ctx context.Context,
	input dto.CreateCityInput,
) error {
	args := m.Called(ctx, input)

	return args.Error(0)
}

func (m *MockCityRepository) Add2User(
	ctx context.Context,
	userID int64,
	input dto.AddCityInput,
) error {
	args := m.Called(ctx, userID, input)

	return args.Error(0)
}

func (m *MockCityRepository) ListCitiesOfUser(
	ctx context.Context,
	userID int64,
	filter dto.ListCitiesFilter,
) ([]model.City, error) {
	args := m.Called(ctx, userID, filter)

	cities, _ := args.Get(0).([]model.City)
	return cities, args.Error(1)
}

func (m *MockCityRepository) GetByName(
	ctx context.Context,
	name string,
) (model.City, error) {
	args := m.Called(ctx, name)

	city, _ := args.Get(0).(model.City)
	return city, args.Error(1)
}

func (m *MockCityRepository) DeleteFromUser(
	ctx context.Context,
	userID, cityID int64,
) error {
	args := m.Called(ctx, userID, cityID)

	return args.Error(0)
}

/// CityServiceTests ///
/// CityServiceTests ///
/// CityServiceTests ///
/// CityServiceTests ///
/// CityServiceTests ///

func TestCityService_Create_Success(t *testing.T) {
	err := config.MustLoad("../../.env")
	require.NoError(t, err)

	client := client.NewOsmClient(config.Cfg.Api)

	repo := new(MockCityRepository)
	cityService := service.NewCityService(repo, client)

	req := dto.CreateCityInput{
		City: "London",
	}

	repo.On("Create", mock.Anything, mock.MatchedBy(func(city dto.CreateCityInput) bool {
		return city.City == "London"
	})).
		Return(nil).
		Once()

	err = cityService.Create(context.Background(), req)

	require.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestCityService_Create_InvalidInput(t *testing.T) {
	err := config.MustLoad("../../.env")
	require.NoError(t, err)

	client := client.NewOsmClient(config.Cfg.Api)

	repo := new(MockCityRepository)
	cityService := service.NewCityService(repo, client)

	req := dto.CreateCityInput{
		City: "",
	}

	err = cityService.Create(context.Background(), req)

	assert.Equal(t, model.ErrInvalidCityInput, err)

	repo.AssertExpectations(t)
}

func TestCityService_Add2User_Success(t *testing.T) {
	err := config.MustLoad("../../.env")
	require.NoError(t, err)

	client := client.NewOsmClient(config.Cfg.Api)

	repo := new(MockCityRepository)
	cityService := service.NewCityService(repo, client)

	req := dto.AddCityInput{
		City: "Paris",
	}

	repo.On("Add2User", mock.Anything, int64(2), mock.MatchedBy(func(city dto.AddCityInput) bool {
		return city.City == "Paris"
	})).
		Return(nil).
		Once()

	err = cityService.Add2User(context.Background(), 2, req)

	require.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestCityService_Add2User_InvalidInput(t *testing.T) {
	err := config.MustLoad("../../.env")
	require.NoError(t, err)

	client := client.NewOsmClient(config.Cfg.Api)

	repo := new(MockCityRepository)
	cityService := service.NewCityService(repo, client)

	req := dto.AddCityInput{
		City: "",
	}

	err = cityService.Add2User(context.Background(), 2, req)

	assert.Equal(t, model.ErrInvalidCityInput, err)

	repo.AssertExpectations(t)
}

func TestCityService_GetByName_Success(t *testing.T) {
	err := config.MustLoad("../../.env")
	require.NoError(t, err)

	client := client.NewOsmClient(config.Cfg.Api)

	repo := new(MockCityRepository)
	userService := service.NewCityService(repo, client)

	expectedCity := model.City{
		CityID: 1,
		City:   "Paris",
	}

	repo.On("GetByName", mock.Anything, "Paris").
		Return(expectedCity, nil).
		Once()

	actualCity, err := userService.GetByName(context.Background(), expectedCity.City)

	require.NoError(t, err)
	assert.Equal(t, expectedCity, actualCity)

	repo.AssertExpectations(t)
}

func TestUserService_GetByName_NotFound(t *testing.T) {
	err := config.MustLoad("../../.env")
	require.NoError(t, err)

	client := client.NewOsmClient(config.Cfg.Api)

	repo := new(MockCityRepository)
	userService := service.NewCityService(repo, client)

	repo.On("GetByName", mock.Anything, "London").
		Return(nil, model.ErrCityNotFound).
		Once()

	_, err = userService.GetByName(context.Background(), "London")

	require.ErrorIs(t, err, model.ErrCityNotFound)

	repo.AssertExpectations(t)
}

func TestCityService_DeleteFromUser_Success(t *testing.T) {
	err := config.MustLoad("../../.env")
	require.NoError(t, err)

	client := client.NewOsmClient(config.Cfg.Api)

	repo := new(MockCityRepository)
	cityService := service.NewCityService(repo, client)

	repo.On("DeleteFromUser", mock.Anything, int64(2), int64(1)).
		Return(nil).
		Once()

	err = cityService.DeleteFromUser(context.Background(), 2, 1)

	require.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestCityService_DeleteFromUser_InvalidInput(t *testing.T) {
	err := config.MustLoad("../../.env")
	require.NoError(t, err)

	client := client.NewOsmClient(config.Cfg.Api)

	repo := new(MockCityRepository)
	cityService := service.NewCityService(repo, client)

	tests := []struct {
		name    string
		userID  int64
		cityID  int64
		wantErr error
	}{
		{
			name:    "Invalid user_id",
			userID:  -2,
			cityID:  1,
			wantErr: model.ErrInvalidUserID,
		},
		{
			name:    "Invalid city_id",
			userID:  2,
			cityID:  -1,
			wantErr: model.ErrInvalidCityID,
		},
		{
			name:    "Invalid user_id and city_id",
			userID:  -2,
			cityID:  -1,
			wantErr: model.ErrInvalidUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err = cityService.DeleteFromUser(
				context.Background(), tt.userID, tt.cityID)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}

	// repo.AssertExpectations(t)
}
