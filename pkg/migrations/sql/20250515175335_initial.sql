-- +goose Up
-- +goose StatementBegin
CREATE TABLE labels (
    id SERIAL PRIMARY KEY,
    key varchar(30) NOT NULL,
    value varchar(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT(now() AT TIME ZONE 'UTC'),
    date TIMESTAMP NOT NULL,
    transaction_type TEXT NOT NULL,
    description TEXT NOT NULL,
    recipient TEXT NOT NULL,
    amount NUMERIC(7,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions_labels (
    transaction_id SERIAL REFERENCES transactions(id) ON DELETE CASCADE,
    label_id SERIAL REFERENCES labels(id) ON DELETE CASCADE,
    CONSTRAINT transaction_label_pkey PRIMARY KEY (transaction_id, label_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions_labels;
DROP TABLE labels;
DROP TABLE transactions;
-- +goose StatementEnd
