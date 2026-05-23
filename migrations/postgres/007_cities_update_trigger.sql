-- +goose Up
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_cities_updated_at ON cities;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_cities_updated_at
BEFORE UPDATE ON cities
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
-- +goose StatementEnd

-- +goose Down
DROP TRIGGER IF EXISTS trg_cities_updated_at ON cities;