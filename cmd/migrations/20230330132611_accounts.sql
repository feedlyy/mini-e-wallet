-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
    id text NOT NULL,
    customer_xid text,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table accounts;
-- +goose StatementEnd
