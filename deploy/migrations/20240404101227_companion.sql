-- +goose Up
-- +goose StatementBegin
CREATE TABLE companions_places(
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    place_id INTEGER,
    date_from DATE,
    date_to DATE
);

CREATE TABLE companions_routes(
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    route_id INTEGER,
    date_from DATE,
    date_to DATE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS companions_places, companions_routes CASCADE;
-- +goose StatementEnd
