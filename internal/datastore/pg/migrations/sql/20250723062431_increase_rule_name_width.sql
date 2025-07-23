-- +goose Up
-- +goose StatementBegin
ALTER TABLE rules ALTER COLUMN id TYPE VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
