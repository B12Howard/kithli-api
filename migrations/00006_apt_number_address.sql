-- +goose Up
-- +goose StatementBegin

ALTER TABLE addresses
    ADD COLUMN apt_number VARCHAR(10)


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE addresses
    DROP COLUMN IF EXISTS apt_number


-- +goose StatementEnd
