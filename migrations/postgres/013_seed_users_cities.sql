-- +goose Up
INSERT INTO users_cities (user_id, city_id)
VALUES
	(2, 1)
ON CONFLICT DO NOTHING;

-- +goose Down
DELETE FROM users_cities 
WHERE user_id = 2 AND city_id = 1;