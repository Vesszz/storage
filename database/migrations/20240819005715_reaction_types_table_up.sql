CREATE TABLE IF NOT EXISTS reaction_types (
    id SERIAL PRIMARY KEY,
    type VARCHAR(31)
);

INSERT INTO reaction_types (type) VALUES ('like'), ('dislike');