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
