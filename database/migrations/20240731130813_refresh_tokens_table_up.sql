CREATE TABLE IF NOT EXISTS refresh_tokens (
    user_id INT,
    fingerprint VARCHAR(255) NOT NULL,
    key UUID UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, fingerprint)
);