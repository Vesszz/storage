CREATE TABLE IF NOT EXISTS files (
    id SERIAL PRIMARY KEY,
    user_id INT,
    key UUID unique,
    path VARCHAR(255) unique,
    time_created TIMESTAMP,
    name VARCHAR(255),
    description VARCHAR(255),
    times_viewed BIGINT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);