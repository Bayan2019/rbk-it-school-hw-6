-- +goose Up
UPDATE users
SET deleted_at = NOW()
WHERE email = 'ayan@example.com';

-- +goose Down
UPDATE users
SET deleted_at = NULL
WHERE email = 'ayan@example.com';