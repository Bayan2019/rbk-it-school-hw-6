package app

import "github.com/jmoiron/sqlx"

func MigrateUp(db *sqlx.DB) error {

	queries := []string{
		"CREATE TYPE IF NOT EXISTS roles AS ENUM ('admin', 'user');",
		`
	CREATE TABLE IF NOT EXISTS users (
		id              BIGSERIAL PRIMARY KEY,
		email           VARCHAR(255) NOT NULL,
		password_hash   VARCHAR(255) NOT NULL,
		first_name      VARCHAR(100) NOT NULL,
		last_name       VARCHAR(100) NOT NULL,
		role            roles NOT NULL DEFAULT 'user',
		is_active       BOOLEAN NOT NULL DEFAULT TRUE,
		created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at      TIMESTAMP NOT NULL DEFAULT NOW(),
		deleted_at      TIMESTAMP
	);`,
		`
	INSERT INTO users (email, password_hash, first_name, last_name, is_active, role)
	VALUES
    ('admin@example.com', '$2a$12$19E8eNXYZEZGelkKqNBIPuGIqptOh4lZsjQHi.Y2V2vZV4V8dFFJe',
        'Admin', 'Role', TRUE, 'admin'),
    ('user@example.com', '$2a$12$wb9gOgmR4OwdVfl2DmfsY.AIK53qR33nYdkPuMAWT8HsIXm7dEpBS', 
        'User', 'Role', TRUE, 'user')
	ON CONFLICT DO NOTHING;`,
		`
	CREATE TABLE cities(
		city_id BIGSERIAL PRIMARY KEY,
		created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		city VARCHAR(255) NOT NULL UNIQUE,
		lat DOUBLE PRECISION NOT NULL,
		lon DOUBLE PRECISION NOT NULL,
		UNIQUE(lat, lon)
	);`,
		`
	CREATE TABLE IF NOT EXISTS users_cities(
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		city_id BIGINT NOT NULL REFERENCES cities(city_id) ON DELETE CASCADE,
		added_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
		deleted_at TIMESTAMP,
		UNIQUE(user_id, city_id)
	);
	`,
		`
	CREATE TABLE IF NOT EXISTS weather_history(
		id BIGSERIAL PRIMARY KEY,
		user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		city VARCHAR(255) NOT NULL REFERENCES cities(city) ON DELETE CASCADE,
		temperature NUMERIC(10, 2) NOT NULL,
		description VARCHAR(255) NOT NULL,
		requested_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func MigrateDown(db *sqlx.DB) error {
	queries := []string{
		"DROP TABLE IF EXISTS weather_history;",
		"DROP TABLE IF EXISTS users_cities;",
		"DROP TABLE IF EXISTS cities;",
		"DELETE FROM users WHERE email IN ('admin@example.com', 'user@example.com');",
		"DROP TABLE IF EXISTS users;",
		"DROP TYPE IF EXISTS roles;",
	}

	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}
