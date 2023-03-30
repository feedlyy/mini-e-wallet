-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id string NOT NULL,
    status string,
    type string,
    amount int,
    reference_id string,
    transaction_at timestamp,
    transaction_by text,
    PRIMARY KEY(id),
    FOREIGN KEY(owned_by) REFERENCES accounts(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table transactions;
-- +goose StatementEnd
