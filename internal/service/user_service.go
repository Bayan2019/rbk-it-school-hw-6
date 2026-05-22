package service

import (
	"context"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
)

type userRepository interface {
	Create(ctx context.Context, input dto.CreateUserInput) (model.User, error)
	List(ctx context.Context, filter dto.ListUsersFilter) ([]model.User, error)
	GetByEmail(ctx context.Context, email string, includeDeleted bool) (model.User, error)
	GetByID(ctx context.Context, id int64, includeDeleted bool) (model.User, error)
	Update(ctx context.Context, id int64, input dto.UpdateUserInput) error
	Delete(ctx context.Context, id int64) error
}

type UserService struct {
	repo userRepository
}

func NewUserService(repo userRepository) *UserService {
	return &UserService{repo: repo}
}

////// methods
////// methods
////// methods

func (s *UserService) Create(
	ctx context.Context,
	input dto.RegisterUserInput,
) (model.User, error) {
	err := input.NormalizeAndValidate()
	if err != nil {
		return model.User{}, err
	}
	cui, err := dto.RegisterUserInput2CreateUserInput(input)
	if err != nil {
		return model.User{}, err
	}
	return s.repo.Create(ctx, cui)
}

func (s *UserService) List(
	ctx context.Context,
	filter dto.ListUsersFilter,
) ([]model.User, error) {
	filter.Normalize()
	return s.repo.List(ctx, filter)
}

func (s *UserService) GetByEmail(
	ctx context.Context,
	email string,
	includeDeleted bool,
) (model.User, error) {
	// if id <= 0 {
	// 	return domain.User{}, domain.ErrInvalidUserID
	// }
	return s.repo.GetByEmail(ctx, email, includeDeleted)
}

func (s *UserService) GetByID(
	ctx context.Context,
	id int64,
	includeDeleted bool,
) (model.User, error) {
	if id <= 0 {
		return model.User{}, model.ErrInvalidUserID
	}
	return s.repo.GetByID(ctx, id, includeDeleted)
}

func (s *UserService) Update(
	ctx context.Context,
	id int64,
	input dto.UpdateUserInput,
) error {
	if id <= 0 {
		return model.ErrInvalidUserID
	}
	if err := input.NormalizeAndValidate(); err != nil {
		return err
	}
	return s.repo.Update(ctx, id, input)
}

func (s *UserService) Delete(
	ctx context.Context,
	id int64,
) error {
	if id <= 0 {
		return model.ErrInvalidUserID
	}
	return s.repo.Delete(ctx, id)
}
