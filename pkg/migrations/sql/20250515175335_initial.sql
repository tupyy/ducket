-- +goose Up
-- +goose StatementBegin
CREATE TABLE labels (
    id SERIAL PRIMARY KEY,
    key varchar(30) NOT NULL,
    value varchar(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(100) PRIMARY KEY,
    date TIMESTAMP NOT NULL,
    transaction_type TEXT NOT NULL,
    description TEXT NOT NULL,
    amount NUMERIC(7,2) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions_labels (
    transaction_id VARCHAR(100) REFERENCES transactions(id) ON DELETE CASCADE,
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
