package model

import (
	"errors"
	"time"
)

var (
	ErrCityNotFound          = errors.New("city not found")
	ErrInvalidCityID         = errors.New("invalid city id")
	ErrInvalidCityInput      = errors.New("invalid city input")
	ErrCityAlreadyAdded2User = errors.New("city already added to user")
	ErrCityNameAlreadyTaken  = errors.New("city's name is already taken")
)

type City struct {
	CityID    int64     `db:"city_id" json:"city_id,omitempty"`
	City      string    `db:"city" json:"city"`
	Lat       float64   `db:"lat" json:"lat"`
	Lon       float64   `db:"lon" json:"lon"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
