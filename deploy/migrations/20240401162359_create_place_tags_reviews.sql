-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE
);

CREATE TABLE IF NOT EXISTS city (
    id SERIAL PRIMARY KEY,
    name VARCHAR
);

CREATE TABLE IF NOT EXISTS district (
    id SERIAL PRIMARY KEY,
    name VARCHAR UNIQUE,
    city_id INTEGER REFERENCES city(id),
    properties JSONB
);

-- TODO: reviews_average
CREATE TABLE IF NOT EXISTS places (
    id SERIAL PRIMARY KEY,
    city_id INTEGER REFERENCES city(id),
    district_id INTEGER,
    properties JSONB,
    name VARCHAR UNIQUE
);

CREATE TABLE places_tags (
    id SERIAL PRIMARY KEY,
    place_id INTEGER REFERENCES places(id) ON DELETE CASCADE,
    tag_id INTEGER REFERENCES tags(id) ON DELETE CASCADE,
    UNIQUE (tag_id, place_id)
);

CREATE TABLE IF NOT EXISTS places_reviews (
    id SERIAL PRIMARY KEY,
    place_id INTEGER,
    author_id INTEGER,
    properties JSONB
);

CREATE TABLE IF NOT EXISTS route_reviews (
    id SERIAL PRIMARY KEY,
    route_id INTEGER,
    author_id INTEGER,
    properties JSONB
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tags, places, places_tags, reviews, city, district CASCADE
-- +goose StatementEnd
