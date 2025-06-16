-- +goose Up
-- +goose StatementBegin
ALTER TABLE rules ADD COLUMN created_at TIMESTAMP DEFAULT now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE rules DROP COLUMN created_at;
-- +goose StatementEnd
