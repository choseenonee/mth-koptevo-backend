-- +goose Up
-- +goose StatementBegin
ALTER TABLE places_reviews
    ADD COLUMN timestamp timestamp;

ALTER TABLE route_reviews
    ADD COLUMN timestamp timestamp;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE places_reviews
    DROP COLUMN timestamp;

ALTER TABLE route_reviews
    DROP COLUMN timestamp;
-- +goose StatementEnd
