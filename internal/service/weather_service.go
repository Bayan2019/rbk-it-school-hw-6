package service

import (
	"context"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
)

type weatherProvider interface {
	GetCurrentWeather(ctx context.Context, lat, lon float64) (model.Weather, error)
}

type weatherRepository interface {
	DoesUserHaveCity(
		ctx context.Context,
		userID int64,
		city string,
	) (bool, error)
	Create(
		ctx context.Context,
		userID int64,
		cityWeather dto.CityWeatherInput,
	) error
	WeatherHistoryOfUser(
		ctx context.Context,
		userID int64,
		filter dto.WeatherHistoryFilter,
	) ([]model.WeatherHistory, error)
}

type WeatherService struct {
	repo     weatherRepository
	provider weatherProvider
}

func NewWeatherService(repo weatherRepository, provider weatherProvider) *WeatherService {
	return &WeatherService{
		repo:     repo,
		provider: provider,
	}
}

////// methods
////// methods
////// methods

func (s *WeatherService) Create(ctx context.Context,
	userID int64,
	city model.City,
) (model.Weather, error) {

	exists, err := s.repo.DoesUserHaveCity(ctx, userID, city.City)
	if err != nil {
		// h.handleError(w, err)
		return model.Weather{}, err
	}

	if !exists {
		return model.Weather{}, model.ErrCityNotFound
	}

	weather, err := s.provider.GetCurrentWeather(ctx, city.Lat, city.Lon)
	if err != nil {
		// h.handleError(w, err)
		return weather, err
	}

	err = s.repo.Create(ctx, userID, dto.CityWeatherInput{
		City:        city.City,
		Temperature: weather.Temperature,
		Description: weather.Description,
	})
	if err != nil {

		return weather, err
	}

	return weather, nil
}

func (s *WeatherService) WeatherHistoryOfUser(
	ctx context.Context,
	userID int64,
	filter dto.WeatherHistoryFilter,
) ([]model.WeatherHistory, error) {

	filter.Normalize()
	return s.repo.WeatherHistoryOfUser(ctx, userID, filter)
}
