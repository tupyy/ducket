BEGIN;

CREATE TABLE IF NOT EXISTS label (
    id SERIAL PRIMARY KEY,
    key varchar(30) NOT NULL,
    value varchar(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS transaction (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT(now() AT TIME ZONE 'UTC'),
    transaction_type TEXT NOT NULL,
    name TEXT NOT NULL,
    user TEXT NOT NULL,
    destination TEXT NOT NULL,
    debit NUMERIC,
    credit NUMERIC
);

CREATE TABLE IF NOT EXISTS saving (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT(now() AT TIME ZONE 'UTC'),
    name TEXT NOT NULL,
    debit NUMERIC,
    credit NUMERIC
);

CREATE TABLE IF NOT EXISTS transactions_labels (
    transaction_id SERIAL REFERENCES transaction(id) ON DELETE CASCADE,
    label_id SERIAL REFERENCES label(id) ON DELETE CASCADE,
    CONSTRAINT transaction_label_pkey PRIMARY KEY (transaction_id, label_id)
);

CREATE TABLE IF NOT EXISTS savings_labels (
    saving_id SERIAL REFERENCES saving(id) ON DELETE CASCADE,
    label_id SERIAL REFERENCES label(id) ON DELETE CASCADE,
    CONSTRAINT saving_label_pkey PRIMARY KEY (saving_id, label_id)
);

END;
