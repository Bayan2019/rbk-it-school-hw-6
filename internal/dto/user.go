package dto

import (
	"strings"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/auth"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
)

type RegisterUserInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  *bool  `json:"is_active,omitempty"`
}

type CreateUserInput struct {
	Email        string `json:"email"`
	PasswordHash string `db:"password_hash"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	IsActive     *bool  `json:"is_active,omitempty"`
}

type UpdateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  *bool  `json:"is_active,omitempty"`
}

type UpdateUserInput struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	IsActive     *bool  `json:"is_active,omitempty"`
}

type ListUsersFilter struct {
	Limit          int    `json:"limit"`
	Offset         int    `json:"offset"`
	Query          string `json:"query"`
	IncludeDeleted bool   `json:"include_deleted"`
}

type UsersResponse struct {
	Data []model.User `json:"data"`
}

type UserResponse struct {
	Data model.User `json:"data"`
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

func (in *RegisterUserInput) NormalizeAndValidate() error {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.Password = strings.TrimSpace(in.Password)
	in.FirstName = strings.TrimSpace(in.FirstName)
	in.LastName = strings.TrimSpace(in.LastName)

	if in.Email == "" || !strings.Contains(in.Email, "@") {
		return model.ErrInvalidUserInput
	}
	if in.Password == "" || in.FirstName == "" || in.LastName == "" {
		return model.ErrInvalidUserInput
	}
	return nil
}

func (in *UpdateUserInput) NormalizeAndValidate() error {
	in.Email = strings.TrimSpace(strings.ToLower(in.Email))
	in.PasswordHash = strings.TrimSpace(in.PasswordHash)
	in.FirstName = strings.TrimSpace(in.FirstName)
	in.LastName = strings.TrimSpace(in.LastName)

	if in.Email == "" || !strings.Contains(in.Email, "@") {
		return model.ErrInvalidUserInput
	}
	if in.PasswordHash == "" || in.FirstName == "" || in.LastName == "" {
		return model.ErrInvalidUserInput
	}
	return nil
}

func (f *ListUsersFilter) Normalize() {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 20
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	f.Query = strings.TrimSpace(
		strings.ToLower(f.Query))
}

func RegisterUserInput2CreateUserInput(rui RegisterUserInput) (CreateUserInput, error) {
	var create CreateUserInput
	create.Email = rui.Email
	create.FirstName = rui.FirstName
	create.LastName = rui.LastName
	create.IsActive = rui.IsActive
	// 6. Безопасность
	// - bcrypt для паролей
	// - не хранить пароли в plain text
	hashPassword, err := auth.HashPassword(rui.Password)
	if err != nil {
		return create, err
	}
	create.PasswordHash = hashPassword

	return create, nil
}

func UpdateUserRequest2UpdateUserInput(uur UpdateUserRequest) (UpdateUserInput, error) {
	var update UpdateUserInput
	update.Email = uur.Email
	update.FirstName = uur.FirstName
	update.LastName = uur.LastName
	update.IsActive = uur.IsActive
	// 6. Безопасность
	// - bcrypt для паролей
	// - не хранить пароли в plain text
	hashPassword, err := auth.HashPassword(uur.Password)
	if err != nil {
		return update, err
	}
	update.PasswordHash = hashPassword

	return update, nil
}
