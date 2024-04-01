-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS routes (
    id SERIAL PRIMARY KEY,
    city_id INTEGER REFERENCES city(id),
    price INTEGER,
    properties JSONB
);

-- TODO: нахуя
CREATE TABLE routes_places (
    id SERIAL PRIMARY KEY,
    position INTEGER,
    route_id INTEGER REFERENCES routes(id) ON DELETE CASCADE,
    place_id INTEGER REFERENCES places(id) ON DELETE CASCADE,
    UNIQUE (route_id, place_id)
);

CREATE TABLE routes_tags (
    id SERIAL PRIMARY KEY,
    route_id INTEGER REFERENCES routes(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    UNIQUE (tag_id, route_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS routes_tags, routes_places, routes CASCADE;
-- +goose StatementEnd
