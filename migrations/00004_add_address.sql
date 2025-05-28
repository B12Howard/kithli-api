-- +goose Up
-- +goose StatementBegin

-- Create Address table
CREATE TABLE addresses (
    id SERIAL PRIMARY KEY,
    street VARCHAR(255),
    city VARCHAR(100),
    state VARCHAR(50),
    postal_code VARCHAR(20),
    location GEOGRAPHY(Point, 4326)
);

CREATE TABLE member_addresses (
  id SERIAL PRIMARY KEY,
  user_id INT REFERENCES members(id) ON DELETE CASCADE,
  address INT REFERENCES addresses(id) ON DELETE CASCADE,
  is_primary BOOLEAN DEFAULT FALSE
);

CREATE TABLE kith_addresses (
  id SERIAL PRIMARY KEY,
  kith_id INT REFERENCES kiths(id) ON DELETE CASCADE,
  address_id INT REFERENCES addresses(id) ON DELETE CASCADE,
  is_primary BOOLEAN DEFAULT FALSE
);


-- Alter Members table
ALTER TABLE members
    DROP COLUMN IF EXISTS postal_code,
    DROP COLUMN IF EXISTS street_address,
    DROP COLUMN IF EXISTS apt_number,
    DROP COLUMN IF EXISTS city,
    DROP COLUMN IF EXISTS state;
-- Alter Kiths table
ALTER TABLE kiths
    DROP COLUMN IF EXISTS postal_code,
    DROP COLUMN IF EXISTS street_address,
    DROP COLUMN IF EXISTS apt_number,
    DROP COLUMN IF EXISTS city,
    DROP COLUMN IF EXISTS state;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Revert Kiths table
ALTER TABLE kiths
    DROP COLUMN IF EXISTS address;

ALTER TABLE kiths
    ADD COLUMN postal_code VARCHAR(10),
    ADD COLUMN street_address VARCHAR,
    ADD COLUMN apt_number VARCHAR,
    ADD COLUMN city VARCHAR,
    ADD COLUMN state VARCHAR(15);

-- Revert Members table
ALTER TABLE members
    DROP COLUMN IF EXISTS address;

ALTER TABLE members
    ADD COLUMN postal_code VARCHAR,
    ADD COLUMN street_address VARCHAR,
    ADD COLUMN apt_number VARCHAR,
    ADD COLUMN city VARCHAR,
    ADD COLUMN state VARCHAR(15);

-- Drop Address table
DROP TABLE IF EXISTS address;

-- +goose StatementEnd
