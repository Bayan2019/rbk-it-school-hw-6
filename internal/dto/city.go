package dto

import (
	"strings"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CityItem struct {
	City string  `db:"city" json:"city"`
	Lat  float64 `db:"lat" json:"lat"`
	Lon  float64 `db:"lon" json:"lon"`
}

type AddCityInput struct {
	City string `json:"city"`
}

type CreateCityInput struct {
	City string  `json:"city"`
	Lat  float64 `json:"lat,omitempty"`
	Lon  float64 `json:"lon,omitempty"`
}

type ListCitiesFilter struct {
	Offset         int  `json:"offset"`
	IncludeDeleted bool `json:"include_deleted"`
}

type CitiesResponse struct {
	Data []model.City `json:"data"`
}

type CityResponse struct {
	Data model.City `json:"data"`
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

func (in *CreateCityInput) NormalizeAndValidate() error {
	in.City = strings.TrimSpace(cases.Title(language.Und).String(in.City))

	if in.City == "" {
		return model.ErrInvalidCityInput
	}

	return nil
}

func (in *AddCityInput) NormalizeAndValidate() error {
	in.City = strings.TrimSpace(cases.Title(language.Und).String(in.City))

	if in.City == "" {
		return model.ErrInvalidCityInput
	}

	return nil
}

func (f *ListCitiesFilter) Normalize() {
	if f.Offset < 0 {
		f.Offset = 0
	}
}
