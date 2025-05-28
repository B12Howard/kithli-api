-- +goose Up
-- +goose StatementBegin

-- Remove availability foreign keys from Kiths table
ALTER TABLE kiths
    DROP COLUMN IF EXISTS days_available,
    DROP COLUMN IF EXISTS times_available;

-- Add kith_id to kith_available (1:1 ownership)
ALTER TABLE kith_available
    ADD COLUMN kith_id INT UNIQUE REFERENCES kiths(id) ON DELETE CASCADE;

-- Add kith_id to times_available (1:1 ownership)
ALTER TABLE times_available
    ADD COLUMN kith_id INT UNIQUE REFERENCES kiths(id) ON DELETE CASCADE;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove kith_id from kith_available and times_available
ALTER TABLE kith_available
    DROP COLUMN IF EXISTS kith_id;

ALTER TABLE times_available
    DROP COLUMN IF EXISTS kith_id;

-- Restore days_available and times_available to kiths
ALTER TABLE kiths
    ADD COLUMN days_available INT REFERENCES kith_available(id),
    ADD COLUMN times_available INT REFERENCES times_available(id);

-- +goose StatementEnd
