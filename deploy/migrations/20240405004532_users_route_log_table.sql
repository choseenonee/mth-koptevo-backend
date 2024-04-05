-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_route_logs (
    user_id INTEGER REFERENCES users(id),
    route_id INTEGER REFERENCES routes(id),
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    UNIQUE (user_id, route_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_route_logs;
-- +goose StatementEnd
