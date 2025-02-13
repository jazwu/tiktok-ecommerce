CREATE TABLE IF NOT EXISTS payment_intents (
    id CHAR(36) PRIMARY KEY,
    amount INT NOT NULL,
    currency CHAR(3) NOT NULL,
    status VARCHAR(20) NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    created_at DATETIME(6) NOT NULL
);

CREATE TABLE IF NOT EXISTS idempotency_keys (
    idempotency_key VARCHAR(255) PRIMARY KEY,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);

CREATE INDEX idx_payment_intents_status ON payment_intents(status);
