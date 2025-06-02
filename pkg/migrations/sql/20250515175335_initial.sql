-- +goose Up
-- +goose StatementBegin
CREATE TABLE rules (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    pattern TEXT NOT NULL
);

CREATE TABLE transactions (
    id VARCHAR(100) PRIMARY KEY,
    date TIMESTAMP NOT NULL,
    transaction_type VARCHAR(30) NOT NULL,
    description TEXT NOT NULL,
    amount NUMERIC(7,2) NOT NULL
);

CREATE TABLE tags (
    value VARCHAR(50) PRIMARY KEY,
    rule_id VARCHAR(100) NOT NULL REFERENCES rules(id) ON DELETE CASCADE
);

CREATE TABLE transactions_tags (
    transaction_id VARCHAR(100) REFERENCES transactions(id) ON DELETE CASCADE,
    tag VARCHAR(50) REFERENCES tags(value) ON DELETE CASCADE,
    CONSTRAINT transaction_tag_pkey PRIMARY KEY (transaction_id, tag)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tags;
DROP TABLE rules;
DROP TABLE transactions_tags;
DROP TABLE transactions;
-- +goose StatementEnd
