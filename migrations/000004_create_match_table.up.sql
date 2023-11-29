CREATE TABLE IF NOT EXISTS matches (
    id int NOT NULL AUTO_INCREMENT,
    category VARCHAR(100) NOT NULL,
    name TEXT NOT NULL,
    home_team_id INT NOT NULL,
    away_team_id INT NOT NULL,
    status int NOT NULL DEFAULT 0,
    close_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    PRIMARY KEY (id)
);