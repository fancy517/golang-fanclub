CREATE TABLE IF NOT EXISTS wallets (
    id int NOT NULL AUTO_INCREMENT,
    user_id int NOT NULL,
    credit double NOT NULL DEFAULT 0,
    deposit_address VARCHAR(255) NOT NULL,
    private_key VARCHAR(255) NOT NULL,
    last_deposit BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id),
    UNIQUE (user_id, deposit_address, private_key)
);
