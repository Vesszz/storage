package dbloader

import (
	"fmt"
	"main/internal/models"
)

func (l *DBLoader) Like(like models.Like) error {
	_, err := l.db.Exec("INSERT INTO likes(user_id, file_id) VALUES ($1, $2)", like.UserID, like.FileID)
	if err != nil {
		return fmt.Errorf("inserting like: %w", err)
	}
	return nil
}

func (l *DBLoader) Unlike(like models.Like) error {
	_, err := l.db.Exec("DELETE FROM likes WHERE user_id = $1 AND file_id = $2", like.UserID, like.FileID)
	if err != nil {
		return fmt.Errorf("deleting like: %w", err)
	}
	return nil
}
