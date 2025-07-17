-- +goose Up
-- +goose StatementBegin
CREATE TABLE rules (
    id VARCHAR(100) PRIMARY KEY,
    created_at TIMESTAMP DEFAULT now(),
    pattern TEXT NOT NULL
);

CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    hash VARCHAR(100) NOT NULL,
    account BIGSERIAL NOT NULL,
    date TIMESTAMP NOT NULL,
    kind VARCHAR(30) NOT NULL,
    content TEXT NOT NULL,
    amount NUMERIC(15,2) NOT NULL
);

CREATE UNIQUE INDEX transaction_hash_idx ON transactions (hash);

CREATE TABLE labels (
    id SERIAL PRIMARY KEY,
    key VARCHAR(255),
    value VARCHAR(100),
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE (key, value)
);

CREATE TABLE rules_labels (
    rule_id VARCHAR(100) REFERENCES rules(id) ON DELETE CASCADE,
    label_id SERIAL REFERENCES labels(id) ON DELETE CASCADE,
    CONSTRAINT rules_labels_pk PRIMARY KEY (rule_id, label_id)
);


CREATE TABLE transactions_labels (
    transaction_id SERIAL REFERENCES transactions(id) ON DELETE CASCADE,
    label_id SERIAL  REFERENCES labels(id) ON DELETE CASCADE,
    rule_id VARCHAR(100) REFERENCES rules(id) ON DELETE CASCADE,
    CONSTRAINT transaction_tag_pk PRIMARY KEY (transaction_id, label_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tags;
DROP TABLE rules;
DROP TABLE rules_labels;
DROP TABLE transactions_labels;
DROP TABLE transactions;
-- +goose StatementEnd
