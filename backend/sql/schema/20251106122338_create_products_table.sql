-- migrations/products.sql

-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    seller_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price INT NOT NULL,
    stock INT NOT NULL DEFAULT 0,
    image_url TEXT,   
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS products_name_idx ON products (name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
