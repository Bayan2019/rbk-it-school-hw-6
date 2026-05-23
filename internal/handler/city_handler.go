package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/middleware"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/go-chi/chi/v5"
)

type cityService interface {
	Create(ctx context.Context, input dto.CreateCityInput) error
	Add2User(ctx context.Context, userID int64, input dto.AddCityInput) error
	ListCitiesOfUser(ctx context.Context, userID int64, filter dto.ListCitiesFilter) ([]model.City, error)
	GetByName(ctx context.Context, name string) (model.City, error)
	DeleteFromUser(ctx context.Context, userID, cityID int64) error
}

type CityHandler struct {
	service cityService
}

func NewCityHandler(service cityService) *CityHandler {
	return &CityHandler{
		service: service,
	}
}

func (h *CityHandler) RegisterAuthRoutes(router chi.Router) {
	router.Post("/cities", h.Add2User)
	router.Get("/cities", h.ListOfUser)
	router.Delete("/cities/{city_id}", h.DeleteFromUser)
}

////// methods
////// methods
////// methods

func (h *CityHandler) Add2User(w http.ResponseWriter, r *http.Request) {

	user, err := middleware.UserFromContext(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	var input dto.AddCityInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "couldn't parse input Body", err)
		return
	}

	city, err := h.service.GetByName(r.Context(), input.City)
	if err != nil {
		if notFound(err) {

			err = h.service.Create(r.Context(), dto.CreateCityInput{
				City: input.City,
			})
			if err != nil {
				h.handleError(w, err)
				return
			}

			err = h.service.Add2User(r.Context(), user.ID, dto.AddCityInput{City: input.City})
			if err != nil {
				h.handleError(w, err)
				return
			}

			middleware.WriteJSON(w, http.StatusCreated, dto.CityResponse{Data: city})
			return
		}
		h.handleError(w, err)
		return
	}

	err = h.service.Add2User(r.Context(), user.ID, dto.AddCityInput{
		City: city.City,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}
	middleware.WriteJSON(w, http.StatusCreated, dto.CityResponse{Data: city})
}

func (h *CityHandler) ListOfUser(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	filter := dto.ListCitiesFilter{
		Offset:         parseIntQuery(r, "offset", 0),
		IncludeDeleted: parseBoolQuery(r, "include_deleted", false),
	}

	cities, err := h.service.ListCitiesOfUser(r.Context(), user.ID, filter)
	if err != nil {
		h.handleError(w, err)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, dto.CitiesResponse{Data: cities})
}

func (h *CityHandler) DeleteFromUser(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.UserFromContext(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	cityID, err := strconv.ParseInt(chi.URLParam(r, "city_id"), 10, 64)
	if err != nil || cityID <= 0 {
		h.handleError(w, model.ErrInvalidCityID)
	}

	if err := h.service.DeleteFromUser(r.Context(), user.ID, cityID); err != nil {
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
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

func (h *CityHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrInvalidCityID), errors.Is(err, model.ErrInvalidCityInput):
		middleware.WriteError(w, http.StatusBadRequest, "invalid city_id of city input", err)
	case errors.Is(err, model.ErrCityNotFound):
		middleware.WriteError(w, http.StatusNotFound, "city is not found", err)
	default:
		// writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "internal server error"})
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", err)
	}
}

func notFound(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "42P01"
	}
	return strings.Contains(strings.ToLower(err.Error()), "not found")
}

func parseIDParam(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, model.ErrInvalidUserID
	}
	return id, nil
}

func parseIntQuery(r *http.Request, key string, fallback int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func parseBoolQuery(r *http.Request, key string, fallback bool) bool {
	value := r.URL.Query().Get(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
