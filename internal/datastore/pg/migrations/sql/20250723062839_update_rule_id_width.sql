-- +goose Up
-- +goose StatementBegin
ALTER TABLE rules_labels ALTER COLUMN rule_id TYPE VARCHAR(255);
ALTER TABLE transactions_labels ALTER COLUMN rule_id TYPE VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
