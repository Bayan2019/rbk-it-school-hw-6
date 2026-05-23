package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/auth"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/handler"
	middle "github.com/Bayan2019/rbk-it-school-hw-6/internal/middleware"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
	"github.com/Bayan2019/rbk-it-school-hw-6/pkg/logger"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

/// UserService ///
/// UserService ///
/// UserService ///
/// UserService ///
/// UserService ///

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Create(
	ctx context.Context,
	input dto.RegisterUserInput,
) (model.User, error) {
	args := m.Called(ctx, input)

	user, _ := args.Get(0).(model.User)
	return user, args.Error(1)
}

func (m *MockUserService) List(
	ctx context.Context,
	filter dto.ListUsersFilter,
) ([]model.User, error) {
	args := m.Called(ctx, filter)

	users, _ := args.Get(0).([]model.User)
	return users, args.Error(1)
}

func (m *MockUserService) GetByID(
	ctx context.Context,
	id int64,
	includeDeleted bool,
) (model.User, error) {
	args := m.Called(ctx, id, includeDeleted)

	user, _ := args.Get(0).(model.User)
	return user, args.Error(1)
}

func (m *MockUserService) GetByEmail(
	ctx context.Context,
	email string,
	includeDeleted bool,
) (model.User, error) {
	args := m.Called(ctx, email, includeDeleted)

	user, _ := args.Get(0).(model.User)
	return user, args.Error(1)
}

func (m *MockUserService) Update(
	ctx context.Context,
	id int64,
	input dto.UpdateUserRequest,
) error {
	args := m.Called(ctx, id, input)

	return args.Error(0)
}

func (m *MockUserService) Delete(
	ctx context.Context,
	id int64,
) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

/// UserHandlerTests ///
/// UserHandlerTests ///
/// UserHandlerTests ///
/// UserHandlerTests ///
/// UserHandlerTests ///

func TestUserHandler_Register_Success(t *testing.T) {
	userService := new(MockUserService)

	logger, close, err := logger.InitializeLogger("")
	defer func() {
		if err := close(); err != nil {
			// and prints any cleanup error to STDERR.
			fmt.Fprintf(os.Stderr, "Failed to close logger: %v\n", err)
		}
	}()
	require.NoError(t, err)

	jwtManager := auth.NewJWTManager([]byte("test-secret"), logger)

	userHandler := handler.NewUserHandler(userService, jwtManager)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	userHandler.RegisterCommonRoutes(r)

	user := model.User{
		ID:        6,
		Email:     "bayan@example.com",
		FirstName: "Bayan",
		LastName:  "User",
		IsActive:  true,
	}

	expectedResponse := dto.UserResponse{
		Data: user,
	}

	trueV := true
	requestBody := dto.RegisterUserInput{
		Email:     "bayan@example.com",
		Password:  "tramp",
		FirstName: "Bayan",
		LastName:  "User",
		IsActive:  &trueV,
	}

	userService.On("Create", mock.Anything, mock.MatchedBy(func(input dto.RegisterUserInput) bool {
		return input.Email == "bayan@example.com"
	})).
		Return(user, nil).
		Once()

	body, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)

	var actualResponse dto.UserResponse
	err = json.Unmarshal(resp.Body.Bytes(), &actualResponse)

	assert.Equal(t, expectedResponse, actualResponse)

	userService.AssertExpectations(t)
}

func TestUserHandler_Register_InvalidJson(t *testing.T) {
	userService := new(MockUserService)

	logger, close, err := logger.InitializeLogger("")
	defer func() {
		if err := close(); err != nil {
			// and prints any cleanup error to STDERR.
			fmt.Fprintf(os.Stderr, "Failed to close logger: %v\n", err)
		}
	}()
	require.NoError(t, err)

	jwtManager := auth.NewJWTManager([]byte("test-secret"), logger)

	userHandler := handler.NewUserHandler(userService, jwtManager)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	userHandler.RegisterCommonRoutes(r)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("{invalid-json"))

	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	userService.AssertNotCalled(t, "Create")
}

func TestUserHandler_Profile_Success(t *testing.T) {
	userService := new(MockUserService)

	logger, close, err := logger.InitializeLogger("")
	defer func() {
		if err := close(); err != nil {
			// and prints any cleanup error to STDERR.
			fmt.Fprintf(os.Stderr, "Failed to close logger: %v\n", err)
		}
	}()
	require.NoError(t, err)

	jwtManager := auth.NewJWTManager([]byte("test-secret"), logger)

	userHandler := handler.NewUserHandler(userService, jwtManager)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Group(func(r chi.Router) {
		r.Use(middle.AuthMiddleware(userHandler.JwtManager))
		userHandler.RegisterAuthRoutes(r)
	})

	user := model.User{
		ID:        2,
		FirstName: "User",
		LastName:  "Role",
		Email:     "user@example.com",
		Role:      auth.RolesUser,
		IsActive:  true,
	}

	expectedResponse := struct {
		Message string
		User    model.User
	}{
		Message: "this is protected profile endpoint",
		User:    user,
	}

	userService.On("GetByEmail", mock.Anything, user.Email, false).
		Return(user, nil).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Content-Type", "application/json")

	token, err := jwtManager.Generate(2, "user@example.com", auth.RolesUser)
	require.NoError(t, err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var actualResponse struct {
		Message string
		User    model.User
	}

	err = json.Unmarshal(resp.Body.Bytes(), &actualResponse)

	assert.Equal(t, expectedResponse, actualResponse)

	userService.AssertExpectations(t)
}

func TestUserHandler_Profile_NoAuthHeader(t *testing.T) {
	userService := new(MockUserService)

	logger, close, err := logger.InitializeLogger("")
	defer func() {
		if err := close(); err != nil {
			// and prints any cleanup error to STDERR.
			fmt.Fprintf(os.Stderr, "Failed to close logger: %v\n", err)
		}
	}()
	require.NoError(t, err)

	jwtManager := auth.NewJWTManager([]byte("test-secret"), logger)

	userHandler := handler.NewUserHandler(userService, jwtManager)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Group(func(r chi.Router) {
		r.Use(middle.AuthMiddleware(userHandler.JwtManager))
		userHandler.RegisterAuthRoutes(r)
	})

	expectedResponse := struct {
		Message string `json:"message"`
		Error   error  `json:"error"`
	}{
		Message: "missing Authorization header",
		Error:   nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)

	var actualResponse struct {
		Message string `json:"message"`
		Error   error  `json:"error"`
	}

	err = json.Unmarshal(resp.Body.Bytes(), &actualResponse)

	assert.Equal(t, expectedResponse, actualResponse)

	userService.AssertNotCalled(t, "GetByEmail")
}

func TestUserHandler_Profile_WrongAuthHeader(t *testing.T) {
	userService := new(MockUserService)

	logger, close, err := logger.InitializeLogger("")
	defer func() {
		if err := close(); err != nil {
			// and prints any cleanup error to STDERR.
			fmt.Fprintf(os.Stderr, "Failed to close logger: %v\n", err)
		}
	}()
	require.NoError(t, err)

	jwtManager := auth.NewJWTManager([]byte("test-secret"), logger)

	userHandler := handler.NewUserHandler(userService, jwtManager)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Group(func(r chi.Router) {
		r.Use(middle.AuthMiddleware(userHandler.JwtManager))
		userHandler.RegisterAuthRoutes(r)
	})

	expectedResponse := struct {
		Message string `json:"message"`
		Error   error  `json:"error"`
	}{
		Message: "invalid Authorization header format",
		Error:   nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Authorization", "Something without Bearer prefix")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotAcceptable, resp.Code)

	var actualResponse struct {
		Message string `json:"message"`
		Error   error  `json:"error"`
	}

	err = json.Unmarshal(resp.Body.Bytes(), &actualResponse)

	assert.Equal(t, expectedResponse, actualResponse)

	userService.AssertNotCalled(t, "GetByEmail")
}

func TestUserHandler_Profile_NoToken(t *testing.T) {
	userService := new(MockUserService)

	logger, close, err := logger.InitializeLogger("")
	defer func() {
		if err := close(); err != nil {
			// and prints any cleanup error to STDERR.
			fmt.Fprintf(os.Stderr, "Failed to close logger: %v\n", err)
		}
	}()
	require.NoError(t, err)

	jwtManager := auth.NewJWTManager([]byte("test-secret"), logger)

	userHandler := handler.NewUserHandler(userService, jwtManager)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Group(func(r chi.Router) {
		r.Use(middle.AuthMiddleware(userHandler.JwtManager))
		userHandler.RegisterAuthRoutes(r)
	})

	expectedResponse := struct {
		Message string `json:"message"`
		Error   error  `json:"error"`
	}{
		Message: "empty token",
		Error:   nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Authorization", "Bearer ")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)

	var actualResponse struct {
		Message string `json:"message"`
		Error   error  `json:"error"`
	}

	err = json.Unmarshal(resp.Body.Bytes(), &actualResponse)

	assert.Equal(t, expectedResponse, actualResponse)

	userService.AssertNotCalled(t, "GetByEmail")
}

func TestUserHandler_Profile_WrongToken(t *testing.T) {
	userService := new(MockUserService)

	logger, close, err := logger.InitializeLogger("")
	defer func() {
		if err := close(); err != nil {
			// and prints any cleanup error to STDERR.
			fmt.Fprintf(os.Stderr, "Failed to close logger: %v\n", err)
		}
	}()
	require.NoError(t, err)

	jwtManager := auth.NewJWTManager([]byte("test-secret"), logger)

	userHandler := handler.NewUserHandler(userService, jwtManager)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))

	r.Group(func(r chi.Router) {
		r.Use(middle.AuthMiddleware(userHandler.JwtManager))
		userHandler.RegisterAuthRoutes(r)
	})

	expectedResponse := struct {
		Message string `json:"message"`
		Error   error  `json:"error"`
	}{
		Message: "invalid token",
		Error:   nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/users/me", nil)
	req.Header.Set("Authorization", "Bearer prefix")
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusUnauthorized, resp.Code)

	var actualResponse struct {
		Message string `json:"message"`
		Error   error  `json:"error"`
	}

	err = json.Unmarshal(resp.Body.Bytes(), &actualResponse)

	assert.Equal(t, expectedResponse, actualResponse)

	userService.AssertNotCalled(t, "GetByEmail")
}
