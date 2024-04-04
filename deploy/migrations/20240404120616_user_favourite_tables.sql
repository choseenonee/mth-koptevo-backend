-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users_favourite_routes (
  user_id INTEGER REFERENCES users(id),
  route_id INTEGER REFERENCES routes(id),
  timestamp TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users_favourite_places (
    user_id INTEGER REFERENCES users(id),
    place_id INTEGER REFERENCES places(id),
    timestamp TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users_favourite_places, users_favourite_routes;
-- +goose StatementEnd
