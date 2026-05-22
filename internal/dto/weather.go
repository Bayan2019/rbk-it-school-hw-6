package dto

import (
	"strings"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type CityWeatherInput struct {
	City        string `db:"city" json:"city"`
	Temperature float64
	Description string
}

type WeatherHistoryFilter struct {
	City   string `json:"city"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type WeatherHistoryResponse struct {
	Data []model.WeatherHistory `json:"data"`
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

func (in *CityWeatherInput) NormalizeAndValidate() error {
	strings.TrimSpace(cases.Title(language.Und).String(in.City))

	if in.City == "" {
		return model.ErrInvalidCityInput
	}

	return nil
}

func (f *WeatherHistoryFilter) Normalize() {
	if f.Limit <= 0 {
		f.Limit = 0
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	f.City = strings.TrimSpace(
		cases.Title(language.Und).String(f.City),
	)
}
