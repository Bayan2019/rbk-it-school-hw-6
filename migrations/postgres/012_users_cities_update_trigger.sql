-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_users_cities_updated_at ON users_cities;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_users_cities_updated_at
BEFORE UPDATE ON users_cities
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS trg_users_cities_updated_at ON users_cities;