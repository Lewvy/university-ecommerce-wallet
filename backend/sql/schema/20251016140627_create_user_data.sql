-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = CURRENT_TIMESTAMP;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE if not exists users (
	id SERIAL PRIMARY KEY,
	name text NOT NULL,
	email text NOT NULL UNIQUE,
	password_hash TEXT,
	google_id TEXT,
	upi_id TEXT UNIQUE,
	phone_number text UNIQUE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	email_verified BOOLEAN NOT NULL DEFAULT FALSE,
	user_type text NOT NULL DEFAULT 'customer',
	version INT NOT NULL DEFAULT 1
);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_timestamp ON users;
DROP TABLE IF EXISTS users;
DROP FUNCTION IF EXISTS trigger_set_timestamp();

-- +goose StatementEnd
