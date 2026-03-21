CREATE SEQUENCE IF NOT EXISTS rules_id_seq START 1;
CREATE SEQUENCE IF NOT EXISTS transactions_id_seq START 1;

CREATE TABLE IF NOT EXISTS rules (
    id          INTEGER PRIMARY KEY DEFAULT nextval('rules_id_seq'),
    name        VARCHAR NOT NULL UNIQUE,
    filter      VARCHAR NOT NULL,
    tags        VARCHAR[] NOT NULL,
    created_at  TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS transactions (
    id          INTEGER PRIMARY KEY DEFAULT nextval('transactions_id_seq'),
    hash        VARCHAR NOT NULL UNIQUE,
    date        DATE NOT NULL,
    account     BIGINT NOT NULL,
    kind        VARCHAR NOT NULL,
    amount      DECIMAL(12,2) NOT NULL,
    content     VARCHAR NOT NULL,
    info        VARCHAR,
    recipient   VARCHAR,
    created_at  TIMESTAMP DEFAULT current_timestamp
);
