-- +goose Up
-- +goose StatementBegin
ALTER TABLE users_place_checkin
    ADD CONSTRAINT unique_user_place UNIQUE(user_id, place_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users_place_checkin
    DROP CONSTRAINT unique_user_place;
-- +goose StatementEnd
