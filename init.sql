CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE wallets
(
    id       UUID            DEFAULT uuid_generate_v4() PRIMARY KEY,
    name     TEXT   NOT NULL UNIQUE,
    currency TEXT   NOT NULL,
    balance  BIGINT NOT NULL DEFAULT 0 CHECK (balance >= 0)
);

CREATE TABLE transactions
(
    id                       UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    operation_type           TEXT   NOT NULL,
    amount                   BIGINT NOT NULL,
    sender_wallet_id         UUID NULL,
    sender_wallet_balance    BIGINT NULL,
    recipient_wallet_id      UUID   NOT NULL,
    recipient_wallet_balance BIGINT NOT NULL,
    processed_at             TIMESTAMP WITHOUT TIME ZONE NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    FOREIGN KEY (sender_wallet_id) REFERENCES wallets (id),
    FOREIGN KEY (recipient_wallet_id) REFERENCES wallets (id)
);
