-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallets (
    id string NOT NULL,
    owned_by text,
    status string,
    enable_at timestamp,
    balance int,
    disable_at timestamp,
    PRIMARY KEY(id),
    FOREIGN KEY(owned_by) REFERENCES accounts(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table wallets;
-- +goose StatementEnd
