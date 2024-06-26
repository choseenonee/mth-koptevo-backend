-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    place_id INTEGER REFERENCES places(id),
    properties JSONB,
    UNIQUE (user_id, place_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notes CASCADE;
-- +goose StatementEnd
