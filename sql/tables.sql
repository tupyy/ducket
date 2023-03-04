BEGIN;

CREATE TABLE IF NOT EXISTS label (
    key varchar(30) PRIMARY KEY,
    value varchar(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS transaction (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT(now() AT TIME ZONE 'UTC'),
    transaction_type TEXT NOT NULL,
    description TEXT NOT NULL,
    recipient TEXT NOT NULL,
    amount NUMERIC NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions_labels (
    transaction_id SERIAL REFERENCES transaction(id) ON DELETE CASCADE,
    label_key VARCHAR(30) REFERENCES label(key) ON DELETE CASCADE,
    CONSTRAINT transaction_label_pkey PRIMARY KEY (transaction_id, label_key)
);

COMMIT;
