package service_test

import (
	"context"
	"testing"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

/// UserRepository ///
/// UserRepository ///
/// UserRepository ///
/// UserRepository ///
/// UserRepository ///

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(
	ctx context.Context,
	input dto.CreateUserInput,
) (model.User, error) {
	args := m.Called(ctx, input)

	user, _ := args.Get(0).(model.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) List(
	ctx context.Context,
	filter dto.ListUsersFilter,
) ([]model.User, error) {
	args := m.Called(ctx, filter)

	users, _ := args.Get(0).([]model.User)
	return users, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(
	ctx context.Context,
	email string,
	includeDeleted bool,
) (model.User, error) {
	args := m.Called(ctx, email, includeDeleted)

	user, _ := args.Get(0).(model.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) GetByID(
	ctx context.Context,
	id int64,
	includeDeleted bool,
) (model.User, error) {
	args := m.Called(ctx, id, includeDeleted)

	user, _ := args.Get(0).(model.User)
	return user, args.Error(1)
}

func (m *MockUserRepository) Update(
	ctx context.Context,
	id int64,
	input dto.UpdateUserInput,
) error {
	args := m.Called(ctx, id, input)

	return args.Error(0)
}

func (m *MockUserRepository) Delete(
	ctx context.Context,
	id int64,
) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

/// UserServiceTests ///
/// UserServiceTests ///
/// UserServiceTests ///
/// UserServiceTests ///
/// UserServiceTests ///

func TestUserService_CreateUser_Success(t *testing.T) {
	repo := new(MockUserRepository)
	userService := service.NewUserService(repo)

	trueV := true
	req := dto.RegisterUserInput{
		FirstName: "Alex",
		LastName:  "Some",
		Email:     "alex@example.com",
		IsActive:  &trueV,
		Password:  "alex123",
	}

	createdUser := model.User{
		ID:        3,
		FirstName: "Alex",
		LastName:  "Some",
		Email:     "alex@example.com",
		IsActive:  trueV,
	}

	repo.On("Create", mock.Anything, mock.MatchedBy(func(user dto.CreateUserInput) bool {
		return user.FirstName == "Alex" && user.Email == "alex@example.com"
	})).
		Return(createdUser, nil).
		Once()

	user, err := userService.Create(context.Background(), req)

	require.NoError(t, err)
	assert.Equal(t, createdUser.ID, user.ID)
	assert.Equal(t, createdUser.Email, user.Email)

	repo.AssertExpectations(t)
}

func TestUserService_CreateUser_InvalidInput(t *testing.T) {
	repo := new(MockUserRepository)
	userService := service.NewUserService(repo)

	trueV := true
	req := dto.RegisterUserInput{
		FirstName: "Alex",
		LastName:  "Some",
		Email:     "alexexample.com",
		IsActive:  &trueV,
		Password:  "alex123",
	}

	_, err := userService.Create(context.Background(), req)

	// require.NoError(t, err)
	assert.Equal(t, model.ErrInvalidUserInput, err)

	repo.AssertExpectations(t)
}

func TestUserService_GetUsers_Success(t *testing.T) {
	repo := new(MockUserRepository)
	userService := service.NewUserService(repo)

	expectedUsers := []model.User{
		model.User{
			ID:        1,
			FirstName: "Admin",
			Email:     "admin@example.com",
		},
		model.User{
			ID:        2,
			FirstName: "User",
			Email:     "user@example.com",
		},
	}

	repo.On("List", mock.Anything, mock.MatchedBy(func(filter dto.ListUsersFilter) bool {
		return true
	})).
		Return(expectedUsers, nil).
		Once()

	actualUsers, err := userService.List(context.Background(), dto.ListUsersFilter{
		Limit:          10,
		Offset:         0,
		IncludeDeleted: false,
	})

	require.NoError(t, err)
	assert.Equal(t, expectedUsers, actualUsers)

	repo.AssertExpectations(t)
}

func TestUserService_GetByEmail_Success(t *testing.T) {
	repo := new(MockUserRepository)
	userService := service.NewUserService(repo)

	expectedUser := model.User{
		ID:        1,
		FirstName: "Admin",
		Email:     "admin@example.com",
		IsActive:  true,
	}

	repo.On("GetByEmail", mock.Anything, "admin@example.com", false).
		Return(expectedUser, nil).
		Once()

	actualUser, err := userService.GetByEmail(context.Background(), expectedUser.Email, false)

	require.NoError(t, err)
	assert.Equal(t, expectedUser, actualUser)

	repo.AssertExpectations(t)
}

func TestUserService_GetByEmail_NotFound(t *testing.T) {
	repo := new(MockUserRepository)
	userService := service.NewUserService(repo)

	repo.On("GetByEmail", mock.Anything, "some@example.com", false).
		Return(nil, model.ErrUserNotFound).
		Once()

	user, err := userService.GetByEmail(context.Background(), "some@example.com", false)

	require.ErrorIs(t, err, model.ErrUserNotFound)
	assert.Equal(t, user, model.User{})

	repo.AssertExpectations(t)
}

func TestUserService_GetUser_Success(t *testing.T) {
	repo := new(MockUserRepository)
	userService := service.NewUserService(repo)

	expectedUser := model.User{
		ID:        1,
		FirstName: "Admin",
		Email:     "admin@example.com",
	}

	repo.On("GetByID", mock.Anything, int64(1), false).
		Return(expectedUser, nil).
		Once()

	actualUser, err := userService.GetByID(context.Background(), 1, false)

	require.NoError(t, err)
	assert.Equal(t, expectedUser, actualUser)

	repo.AssertExpectations(t)
}

func TestUserService_GetByID_InvalidID(t *testing.T) {
	repo := new(MockUserRepository)
	userService := service.NewUserService(repo)

	user, err := userService.GetByID(context.Background(), 0, false)

	require.ErrorIs(t, err, model.ErrInvalidUserID)
	assert.Equal(t, user, model.User{})

	repo.AssertNotCalled(t, "GetByID")
}

func TestUserService_GetUser_NotFound(t *testing.T) {
	repo := new(MockUserRepository)
	userService := service.NewUserService(repo)

	repo.On("GetByID", mock.Anything, int64(999), false).
		Return(model.User{}, model.ErrUserNotFound).
		Once()

	user, err := userService.GetByID(context.Background(), 999, false)

	require.ErrorIs(t, err, model.ErrUserNotFound)
	assert.Equal(t, user, model.User{})

	repo.AssertExpectations(t)
}

func TestUserService_Update_Success(t *testing.T) {
	repo := new(MockUserRepository)
	userService := service.NewUserService(repo)

	trueV := true
	req := dto.UpdateUserRequest{
		FirstName: "Alex",
		LastName:  "Some",
		Email:     "admin@example.com",
		IsActive:  &trueV,
		Password:  "admin123",
	}

	// input, err := dto.UpdateUserRequest2UpdateUserInput(req)

	repo.On("Update", mock.Anything, int64(1), mock.MatchedBy(func(user dto.UpdateUserInput) bool {
		return user.FirstName == "Alex" && user.Email == "admin@example.com"
	})).
		Return(nil).
		Once()

	err := userService.Update(context.Background(), 1, req)

	require.NoError(t, err)
	// assert.Equal(t, createdUser.ID, user.ID)
	// assert.Equal(t, createdUser.Email, user.Email)

	repo.AssertExpectations(t)
}

func TestUserService_Update_InvalidInput(t *testing.T) {
	repo := new(MockUserRepository)
	userService := service.NewUserService(repo)

	trueV := true
	req := dto.UpdateUserRequest{
		FirstName: "Admin",
		LastName:  "Some",
		Email:     "alexexample.com",
		IsActive:  &trueV,
		Password:  "alex123",
	}

	err := userService.Update(context.Background(), 1, req)

	// require.NoError(t, err)
	assert.Equal(t, model.ErrInvalidUserInput, err)

	repo.AssertExpectations(t)
}

func TestUserService_Update_NotFound(t *testing.T) {
	repo := new(MockUserRepository)
	userService := service.NewUserService(repo)

	trueV := true
	req := dto.UpdateUserRequest{
		FirstName: "Admin",
		LastName:  "Some",
		Email:     "admin@example.com",
		IsActive:  &trueV,
		Password:  "alex123",
	}

	repo.On("Update", mock.Anything, int64(999), mock.MatchedBy(func(user dto.UpdateUserInput) bool {
		return user.FirstName == "Admin" && user.Email == "admin@example.com"
	})).
		Return(model.ErrUserNotFound).
		Once()

	err := userService.Update(context.Background(), 999, req)

	// require.NoError(t, err)
	assert.Equal(t, model.ErrUserNotFound, err)

	repo.AssertExpectations(t)
}
