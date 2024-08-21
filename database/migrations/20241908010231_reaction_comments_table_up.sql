CREATE TABLE IF NOT EXISTS reaction_comments (
    user_id INT,
    comment_id INT,
    type_id INT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,
    FOREIGN KEY (type_id) REFERENCES reaction_types(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, comment_id)
);