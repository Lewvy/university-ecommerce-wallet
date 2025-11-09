-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tokens (
    hash BYTEA PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL,
    scope TEXT NOT NULL
);

-- Create an index on the user_id and scope columns for fast lookups
CREATE INDEX IF NOT EXISTS tokens_user_id_scope_idx ON tokens (user_id, scope);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP table if exists tokens;
-- +goose StatementEnd
