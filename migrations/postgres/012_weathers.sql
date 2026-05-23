-- +goose Up
CREATE TABLE weathers(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    city VARCHAR(255) NOT NULL REFERENCES cities(city) ON DELETE CASCADE,
    temperature NUMERIC(10, 2) NOT NULL,
    description VARCHAR(255) NOT NULL,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, city, requested_at)
);

-- +goose Down
DROP TABLE weathers;