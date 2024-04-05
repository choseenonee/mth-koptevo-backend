-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS trips (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    date_start DATE,
    date_end DATE,
    properties JSONB
);

CREATE TABLE IF NOT EXISTS trip_places (
    trip_id INTEGER REFERENCES trips(id),
    day INTEGER,
    position INTEGER,
    place_id INTEGER REFERENCES places(id)
);

CREATE TABLE IF NOT EXISTS trip_routes (
    trip_id INTEGER REFERENCES trips(id),
    day INTEGER,
    position INTEGER,
    route_id INTEGER REFERENCES routes(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS trip_routes, trip_places, trips;
-- +goose StatementEnd
