-- +goose Up
-- +goose StatementBegin
CREATE TABLE rules (
    id VARCHAR(12) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    pattern TEXT NOT NULL
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    hash VARCHAR(100) NOT NULL,
    date TIMESTAMP NOT NULL,
    kind VARCHAR(30) NOT NULL,
    content TEXT NOT NULL,
    amount NUMERIC(7,2) NOT NULL
);

CREATE UNIQUE INDEX transaction_hash_idx ON transactions (hash);

CREATE TABLE tags (
    value VARCHAR(50) PRIMARY KEY
);

CREATE TABLE rules_tags (
    rule_id VARCHAR(12) REFERENCES rules(id) ON DELETE CASCADE,
    tag VARCHAR(50) REFERENCES tags(value) ON DELETE CASCADE,
    CONSTRAINT rules_tags_pk PRIMARY KEY (rule_id, tag)
);

CREATE TABLE transactions_tags (
    transaction_id SERIAL REFERENCES transactions(id) ON DELETE CASCADE,
    tag_id VARCHAR(50) REFERENCES tags(value) ON DELETE CASCADE,
    rule_id VARCHAR(12) REFERENCES rules(id) ON DELETE CASCADE,
    CONSTRAINT transaction_tag_pk PRIMARY KEY (transaction_id, tag_id, rule_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tags;
DROP TABLE rules;
DROP TABLE rules_tags;
DROP TABLE transactions_tags;
DROP TABLE transactions;
-- +goose StatementEnd
