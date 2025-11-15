-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS wallet_transactions (
    id SERIAL PRIMARY KEY,
    
    user_id INT NOT NULL,
    amount BIGINT NOT NULL,

    transaction_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    transaction_type VARCHAR(30) NOT NULL,

    related_user_id INT,

    razorpay_order_id TEXT,
    razorpay_payment_id TEXT,

    metadata JSONB,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id) 
        REFERENCES users(id)
        ON DELETE SET NULL,
    
    CONSTRAINT fk_related_user
        FOREIGN KEY(related_user_id)
        REFERENCES users(id)
        ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_razorpay_order_id ON wallet_transactions(razorpay_order_id);

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON wallet_transactions
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TRIGGER IF EXISTS set_timestamp ON wallet_transactions;
DROP TABLE IF EXISTS wallet_transactions;

-- +goose StatementEnd
