-- +goose Up
-- +goose StatementBegin
ALTER TABLE transactions ADD COLUMN info TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions DROP COLUMN info;
-- +goose StatementEnd
