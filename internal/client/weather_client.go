package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
)

type WeatherClient struct {
	httpClient *http.Client
	baseURL    string
}

func NewWeatherClient() *WeatherClient {
	return &WeatherClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://api.open-meteo.com/v1/forecast",
	}
}

////// methods
////// methods
////// methods

func (c *WeatherClient) GetCurrentWeather(ctx context.Context, lat, lon float64) (model.Weather, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return model.Weather{}, fmt.Errorf("parse base url: %w", err)
	}

	q := u.Query()
	q.Set("latitude", fmt.Sprintf("%.4f", lat))
	q.Set("longitude", fmt.Sprintf("%.4f", lon))
	q.Set("current_weather", "true")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return model.Weather{}, fmt.Errorf("create request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return model.Weather{}, fmt.Errorf("call external api: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return model.Weather{}, fmt.Errorf("external api returned status: %d", resp.StatusCode)
	}

	var result model.OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return model.Weather{}, fmt.Errorf("decode external api response: %w", err)
	}

	return model.Weather{
		Temperature: result.CurrentWeather.Temperature,
		Description: mapWeatherCode(result.CurrentWeather.Weathercode),
	}, nil
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

func mapWeatherCode(code int) string {
	switch code {
	case 0:
		return "sunny"
	case 1, 2, 3:
		return "cloudy"
	case 45, 48:
		return "fog"
	case 51, 53, 55:
		return "drizzle"
	case 61, 63, 65:
		return "rain"
	case 71, 73, 75:
		return "snow"
	case 95:
		return "storm"
	default:
		return "unknown"
	}
}
