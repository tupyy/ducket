-- +goose Up
-- +goose StatementBegin
ALTER TABLE tags ADD COLUMN created_at TIMESTAMP DEFAULT now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tags DROP COLUMN created_at;
-- +goose StatementEnd
