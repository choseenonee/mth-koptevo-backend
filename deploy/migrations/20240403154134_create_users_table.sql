-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    properties JSONB
);

CREATE TABLE IF NOT EXISTS users_place_checkin (
    user_id INTEGER REFERENCES users(id),
    place_id INTEGER REFERENCES places(id),
    timestamp timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users, users_place_checkin CASCADE;
-- +goose StatementEnd