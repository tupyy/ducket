-- +goose Up
-- +goose StatementBegin
CREATE TABLE tags (
    id SERIAL PRIMARY KEY,
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
    tag_id SERIAL REFERENCES tags(id) ON DELETE CASCADE,
    CONSTRAINT transaction_tag_pkey PRIMARY KEY (transaction_id, tag_id)
);

CREATE TABLE rules (
    id SERIAL PRIMARY KEY,
    tag_id SERIAL ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    pattern TEXT NOT NULL,
    CONSTRAINT fk_tag FOREIGN KEY(tag_id) REFERENCES tags(id) 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE transactions_labels;
DROP TABLE transactions;
DROP TABLE rules;
DROP TABLE tags;
-- +goose StatementEnd
