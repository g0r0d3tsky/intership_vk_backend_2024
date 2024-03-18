-- +goose Up
-- +goose StatementBegin
ALTER TABLE movies
    ALTER COLUMN rating TYPE real;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE movies
    ALTER COLUMN rating TYPE integer;
-- +goose StatementEnd