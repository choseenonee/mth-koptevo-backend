-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN login VARCHAR;

ALTER TABLE users
    ADD COLUMN password VARCHAR;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN login;

ALTER TABLE users
    DROP COLUMN password;
-- +goose StatementEnd
