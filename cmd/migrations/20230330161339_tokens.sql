-- +goose Up
-- +goose StatementBegin
CREATE TABLE tokens (
    id SERIAL NOT NULL,
    account_id text,
    token text,
    expiration timestamp NOT NULL,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(id),
    FOREIGN KEY(account_id) REFERENCES accounts(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table tokens;
-- +goose StatementEnd
