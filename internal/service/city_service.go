package service

import (
	"context"
	"strconv"
	"strings"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type osmProvider interface {
	GetInfoOfCity(ctx context.Context, city string) (model.Place, error)
}

type cityRepository interface {
	Create(
		ctx context.Context,
		input dto.CreateCityInput,
	) error
	Add2User(
		ctx context.Context,
		userID int64,
		input dto.AddCityInput,
	) error
	ListCitiesOfUser(
		ctx context.Context,
		userID int64,
		filter dto.ListCitiesFilter,
	) ([]model.City, error)
	GetByName(
		ctx context.Context,
		name string,
	) (model.City, error)
	DeleteFromUser(
		ctx context.Context,
		userID, cityID int64,
	) error
}

type CityService struct {
	repo     cityRepository
	provider osmProvider
}

func NewCityService(
	repo cityRepository,
	provider osmProvider,
) *CityService {
	return &CityService{
		repo:     repo,
		provider: provider,
	}
}

////// methods
////// methods
////// methods

func (s *CityService) Create(
	ctx context.Context,
	input dto.CreateCityInput,
) error {
	if err := input.NormalizeAndValidate(); err != nil {
		return err
	}

	if input.Lat == 0.0 && input.Lon == 0.0 {
		place, err := s.provider.GetInfoOfCity(ctx, strings.TrimSpace(strings.ToLower(input.City)))
		if err != nil {
			return err
		}
		lat, err := strconv.ParseFloat(place.Lat, 64)
		if err != nil {
			// h.handleError(w, err)
			return err
		}
		lon, err := strconv.ParseFloat(place.Lon, 64)
		if err != nil {
			// h.handleError(w, err)
			return err
		}
		input.Lat = lat
		input.Lon = lon
	}

	return s.repo.Create(ctx, input)
}

func (s *CityService) Add2User(
	ctx context.Context,
	userID int64,
	input dto.AddCityInput,
) error {
	if err := input.NormalizeAndValidate(); err != nil {
		return err
	}
	return s.repo.Add2User(ctx, userID, input)
}

func (s *CityService) ListCitiesOfUser(
	ctx context.Context,
	userID int64,
	filter dto.ListCitiesFilter,
) ([]model.City, error) {
	filter.Normalize()
	return s.repo.ListCitiesOfUser(ctx, userID, filter)
}

func (s *CityService) GetByName(
	ctx context.Context,
	name string,
) (model.City, error) {
	return s.repo.GetByName(ctx, strings.TrimSpace(cases.Title(language.Und).String(name)))
}

func (s *CityService) DeleteFromUser(
	ctx context.Context,
	userID, cityID int64,
) error {
	if userID <= 0 {
		return model.ErrInvalidUserID
	}
	if cityID <= 0 {
		return model.ErrInvalidCityID
	}
	return s.repo.DeleteFromUser(ctx, userID, cityID)
}
