package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/middleware"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
)

type weatherService interface {
	CreateHistory(ctx context.Context, userID int64, city model.City) (model.Weather, error)
	WeatherHistoryOfUser(ctx context.Context, userID int64, filter dto.WeatherHistoryFilter) ([]model.WeatherHistory, error)
}

type WeatherHandler struct {
	CityService    cityService
	WeatherService weatherService
}

func NewWeatherHandler(city cityService, weather weatherService) *WeatherHandler {
	return &WeatherHandler{
		CityService:    city,
		WeatherService: weather,
	}
}

////// methods
////// methods
////// methods

func (h *WeatherHandler) GetWeatherOfUserCities(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}
	cities, err := h.CityService.ListCitiesOfUser(r.Context(), user.ID, dto.ListCitiesFilter{})
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "couldn't get cities of user", err)
		return
	}

	results := []model.WeatherHistory{}

	for _, city := range cities {
		weather, err := h.WeatherService.CreateHistory(r.Context(), user.ID, city)
		if err != nil {
			h.handleError(w, err)
			return
		}
		results = append(results, model.WeatherHistory{
			UserID: user.ID,
			City:   city.City,
			// RequestedAt: weather.RequestedAt,
			Temperature: weather.Temperature,
			Description: weather.Description,
		})
	}

	middleware.WriteJSON(w, http.StatusOK, results)
}

func (h *WeatherHandler) GetWeatherHistoryOfUser(w http.ResponseWriter, r *http.Request) {
	filter := dto.WeatherHistoryFilter{
		Limit:  parseIntQuery(r, "limit", 20),
		Offset: parseIntQuery(r, "offset", 0),
		City:   r.URL.Query().Get("city"),
	}
	user, err := middleware.UserFromContext(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	result, err := h.WeatherService.WeatherHistoryOfUser(r.Context(), user.ID, filter)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "couldn't get weather history", err)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, dto.WeatherHistoryResponse{Data: result})
}

////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions

func (h *WeatherHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrInvalidUserID), errors.Is(err, model.ErrInvalidUserInput):
		middleware.WriteError(w, http.StatusBadRequest, "invalid user_id or invalid user input", err)
	case errors.Is(err, model.ErrCityNotFound):
		middleware.WriteError(w, http.StatusNotFound, "not found", err)
	default:
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", err)
		// writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
	}
}
