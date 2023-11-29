CREATE TABLE IF NOT EXISTS settings (
    deposit_wallet_private_key VARCHAR(255) NOT NULL DEFAULT '',
    last_block_height BIGINT NOT NULL DEFAULT 0
);

INSERT INTO settings (deposit_wallet_private_key, last_block_height) VALUES ('', 0);