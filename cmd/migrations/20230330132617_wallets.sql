-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets (
    id text NOT NULL,
    owned_by text,
    status text,
    enabled_at timestamp,
    balance int,
    disabled_at timestamp,
    PRIMARY KEY(id),
    FOREIGN KEY(owned_by) REFERENCES accounts(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table wallets;
-- +goose StatementEnd
