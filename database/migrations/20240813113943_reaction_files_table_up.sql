CREATE TABLE IF NOT EXISTS reaction_files (
    user_id INT,
    file_id INT,
    type_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
    FOREIGN KEY (type_id) REFERENCES reaction_types(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, file_id)
);