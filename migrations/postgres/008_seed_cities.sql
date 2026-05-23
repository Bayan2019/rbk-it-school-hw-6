-- +goose Up
INSERT INTO cities (city, lat, lon)
VALUES
	('Paris', 48.8534951, 2.3483915),
	('Berlin', 52.5173885, 13.3951309)
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM cities WHERE city IN ('Paris', 'Berlin');