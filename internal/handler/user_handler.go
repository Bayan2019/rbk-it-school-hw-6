package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/auth"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/middleware"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
)

type userService interface {
	Create(ctx context.Context, input dto.RegisterUserInput) (model.User, error)
	List(ctx context.Context, filter dto.ListUsersFilter) ([]model.User, error)
	GetByID(ctx context.Context, id int64, includeDeleted bool) (model.User, error)
	GetByEmail(ctx context.Context, email string, includeDeleted bool) (model.User, error)
	Update(ctx context.Context, id int64, input dto.UpdateUserRequest) error
	Delete(ctx context.Context, id int64) error
}

type UserHandler struct {
	service    userService
	JwtManager *auth.JWTManager
}

/// json
/// json
/// json
/// json
/// json

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

/// json
/// json
/// json
/// json
/// json

func NewUserHandler(service userService, jwtManager *auth.JWTManager) *UserHandler {
	return &UserHandler{
		service:    service,
		JwtManager: jwtManager,
	}
}

////// methods
////// methods
////// methods

// 1. Аутентификация
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	// 1. Аутентификация
	// - регистрация пользователя
	var input dto.RegisterUserInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json body", err)
		return
	}

	user, err := h.service.Create(r.Context(), input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// w.Header().Set("Location", "/api/v1/users/"+strconv.FormatInt(user.ID, 10))
	middleware.WriteJSON(w, http.StatusCreated, dto.UserResponse{Data: user})
}

// 1. Аутентификация
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json body", err)
		return
	}

	user, err := h.service.GetByEmail(r.Context(), req.Email, false)
	if err == model.ErrUserNotFound {
		middleware.WriteError(w, http.StatusUnauthorized, "email not found", err)
		return
	}
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "internal error", err)
		return
	}

	if !auth.CheckPassword(req.Password, user.PasswordHash) {
		middleware.WriteError(w, http.StatusUnauthorized, "invalid email or password", nil)
		return
	}

	token, err := h.JwtManager.Generate(user.ID, user.Email, user.Role)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "token generation failed", err)
		return
	}

	// 1. Аутентификация
	// - возвращает access_token (JWT)
	middleware.WriteJSON(w, http.StatusOK, loginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
	})
}

func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	userCtx, err := middleware.UserFromContext(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	user, err := h.service.GetByEmail(r.Context(), userCtx.Email, false)
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "couldn't get user", err)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "this is protected profile endpoint",
		"user":    user,
	})
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := dto.ListUsersFilter{
		Limit:          parseIntQuery(r, "limit", 20),
		Offset:         parseIntQuery(r, "offset", 0),
		Query:          r.URL.Query().Get("q"),
		IncludeDeleted: parseBoolQuery(r, "include_deleted", false),
	}

	users, err := h.service.List(r.Context(), filter)
	if err != nil {
		h.handleError(w, err)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, dto.UsersResponse{Data: users})
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.handleError(w, err)
		return
	}

	user, err := h.service.GetByID(r.Context(), id, parseBoolQuery(r, "include_deleted", false))
	if err != nil {
		h.handleError(w, err)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, dto.UserResponse{Data: user})
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	userCtx, err := middleware.UserFromContext(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusUnauthorized, "unauthorized", err)
		return
	}

	var input dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid json body", err)
		return
	}

	// update, err := dto.UpdateUserRequest2UpdateUserInput(input)
	// if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
	// 	middleware.WriteError(w, http.StatusInternalServerError, "error mapping UpdateUserRequest to UpdateUserInput", err)
	// 	return
	// }
	err = h.service.Update(r.Context(), userCtx.ID, input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	middleware.WriteJSON(w, http.StatusOK, map[string]string{
		"message": "user is updated",
	})
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDParam(r)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
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

func (h *UserHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrInvalidUserID), errors.Is(err, model.ErrInvalidUserInput):
		middleware.WriteError(w, http.StatusBadRequest, "invalid user_id and user input", err)
	case errors.Is(err, model.ErrUserNotFound):
		middleware.WriteError(w, http.StatusNotFound, "user not found", err)
	case errors.Is(err, model.ErrEmailAlreadyTaken):
		middleware.WriteError(w, http.StatusConflict, "email is already in use", err)
	default:
		middleware.WriteError(w, http.StatusInternalServerError, "internal server error", err)
		// writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
	}
}
