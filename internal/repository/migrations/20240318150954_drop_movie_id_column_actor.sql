-- +goose Up
-- +goose StatementBegin
ALTER TABLE actors DROP COLUMN movie_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE actors ADD COLUMN movie_id uuid;
-- +goose StatementEnd