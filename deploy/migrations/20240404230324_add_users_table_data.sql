-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN login VARCHAR;

ALTER TABLE users
    ADD COLUMN password VARCHAR;

ALTER TABLE users
    ADD CONSTRAINT unique_login UNIQUE (login);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP COLUMN login;

ALTER TABLE users
    DROP COLUMN password;

ALTER TABLE users
    DROP CONSTRAINT unique_login;
-- +goose StatementEnd
