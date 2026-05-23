-- +goose Up
CREATE INDEX IF NOT EXISTS idx_weather_user_id_city 
ON weathers (user_id, city);

-- +goose Down
DROP INDEX IF EXISTS idx_weather_user_id_city;