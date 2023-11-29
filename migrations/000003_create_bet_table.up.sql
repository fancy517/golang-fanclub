CREATE TABLE IF NOT EXISTS bets (
    id int NOT NULL AUTO_INCREMENT,
    user_id int NOT NULL,
    match_id int NOT NULL DEFAULT 0,
    amount double NOT NULL DEFAULT 0,
    dir int NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id)
);