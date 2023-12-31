CREATE TABLE IF NOT EXISTS notifications (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    user_id INT NOT NULL,
    source VARCHAR(20) NOT NULL,
    foreign_key INT NOT NULL,
    fk1 INT NOT NULL DEFAULT 0,
    fk2 INT NOT NULL DEFAULT 0,
    amount DOUBLE NOT NULL DEFAULT 0,
    memo VARCHAR(100) NOT NULL DEFAULT '',
    created_at TIMESTAMP DEFAULT NOW()
);