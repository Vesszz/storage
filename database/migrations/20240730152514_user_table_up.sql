CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    salt CHAR(8) NOT NULL
);