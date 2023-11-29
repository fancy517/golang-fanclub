CREATE TABLE IF NOT EXISTS tickets (
    id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    subject VARCHAR(255) NOT NULL DEFAULT '',
    message TEXT NOT NULL,
    status VARCHAR(100) NOT NULL DEFAULT 'unread',
    created_at TIMESTAMP DEFAULT NOW()
);