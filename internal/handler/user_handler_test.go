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
	resp := httptest.NewRecorder()

	r.ServeHTTP(resp, req)

	require.Equal(t, http.StatusCreated, resp.Code)

	var actualResponse dto.UserResponse
	err = json.Unmarshal(resp.Body.Bytes(), &actualResponse)

	assert.Equal(t, expectedResponse, actualResponse)

	userService.AssertExpectations(t)
}
