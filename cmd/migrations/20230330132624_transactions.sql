-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id text NOT NULL,
    status text,
    type text,
    amount int,
    reference_id text UNIQUE,
    transaction_at timestamp,
    transaction_by text,
    PRIMARY KEY(id),
    FOREIGN KEY(transaction_by) REFERENCES accounts(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table transactions;
-- +goose StatementEnd
