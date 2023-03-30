-- +goose Up
-- +goose StatementBegin
CREATE TABLE accounts (
    id string NOT NULL,
    customer_xid text,
    PRIMARY KEY(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- SELECT 'down SQL query';
drop table accounts;
-- +goose StatementEnd
