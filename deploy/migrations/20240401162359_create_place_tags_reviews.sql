-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE
);

CREATE TABLE IF NOT EXISTS city (
    id SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE
);

CREATE TABLE IF NOT EXISTS district (
    id SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE,
    city_id INTEGER REFERENCES city(id),
    properties JSONB
);

CREATE TABLE IF NOT EXISTS places (
    id SERIAL PRIMARY KEY,
    reviews_avg INTEGER,
    city_id INTEGER REFERENCES city(id),
    district_id INTEGER REFERENCES district(id),
    properties JSONB
);

CREATE TABLE places_tags (
    id SERIAL PRIMARY KEY,
    place_id INTEGER REFERENCES places(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    UNIQUE (tag_id, place_id)
);

-- TODO: сделать референс на users
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    place_id INTEGER REFERENCES places(id),
    author_id INTEGER,
    properties JSONB
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tags, places, places_tags, reviews, city, district CASCADE
-- +goose StatementEnd
