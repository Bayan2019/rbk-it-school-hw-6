package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Bayan2019/rbk-it-school-hw-6/internal/dto"
	"github.com/Bayan2019/rbk-it-school-hw-6/internal/model"
	"github.com/jmoiron/sqlx"
)

type WeatherRepository struct {
	db *sqlx.DB
}

func NewWeatherRepository(db *sqlx.DB) *WeatherRepository {
	return &WeatherRepository{db: db}
}

////// methods
////// methods
////// methods

func (r *WeatherRepository) CreateHistory(
	ctx context.Context,
	userID int64,
	cityWeather dto.CityWeatherInput,
) error {

	query := `
		INSERT INTO weather_history (user_id, city, temperature, description)
		VALUES (:user_id, :city, :temperature, :description)
	`

	args := map[string]any{
		"user_id":     userID,
		"city":        cityWeather.City,
		"temperature": cityWeather.Temperature,
		"description": cityWeather.Description,
	}

	result, err := r.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("result.RowsAffected(): %v", err)
	}
	if rowsAffected == 0 {
		return errors.New("weather of city is not created")
	}

	return nil
}

func (r *WeatherRepository) WeatherHistoryOfUser(
	ctx context.Context,
	userID int64,
	filter dto.WeatherHistoryFilter,
) ([]model.WeatherHistory, error) {
	builder := strings.Builder{}
	builder.WriteString(`
		SELECT user_id, city, temperature, description, requested_at
		FROM weather_history
		WHERE user_id = :user_id
	`)

	args := map[string]any{
		"user_id": userID,
	}

	if filter.City != "" {
		builder.WriteString(" AND city = :city")
		args["city"] = filter.City
	}

	builder.WriteString(" ORDER BY requested_at DESC")

	if filter.Limit != 0 {
		builder.WriteString(" LIMIT :limit")
		args["limit"] = filter.Limit
	}

	if filter.Offset != 0 {
		builder.WriteString(" OFFSET :offset")
		args["offset"] = filter.Offset
	}

	query, queryArgs, err := sqlx.Named(builder.String(), args)
	if err != nil {
		return []model.WeatherHistory{}, errors.New("sqlx.Named")
	}
	query = r.db.Rebind(query)

	var results []model.WeatherHistory
	if err := r.db.SelectContext(ctx, &results, query, queryArgs...); err != nil {
		return []model.WeatherHistory{}, errors.New("r.db.SelectContext")
	}

	if len(results) == 0 {
		return []model.WeatherHistory{}, nil
	}

	return results, nil
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
