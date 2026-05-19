package model

import (
	"errors"
	"time"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/auth"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidUserID    = errors.New("invalid user id")
	ErrInvalidUserInput = errors.New("invalid user input")
	// 1. Аутентификация
	// - email должен быть уникальным
	ErrEmailAlreadyTaken = errors.New("email already exists")
	ErrHashingPassword   = errors.New("error hashing password")
)

type User struct {
	ID           int64      `db:"id" json:"id"`
	Email        string     `db:"email" json:"email"`
	PasswordHash string     `db:"password_hash,omitempty" json:"-"`
	FirstName    string     `db:"first_name" json:"first_name"`
	LastName     string     `db:"last_name" json:"last_name"`
	Role         auth.Roles `db:"role" json:"role"`
	IsActive     bool       `db:"is_active" json:"is_active"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
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
