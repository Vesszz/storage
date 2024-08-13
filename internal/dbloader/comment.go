package dbloader

import (
	"fmt"
	"main/internal/models"
)

func (l *DBLoader) CreateComment(comment *models.Comment) (*models.Comment, error) {
	err := l.db.QueryRow("INSERT INTO comments(user_id, file_id, time_created, text) VALUES ($1, $2, $3, $4) RETURNING id", comment.UserID, comment.FileID, comment.TimeCreated, comment.Text).Scan(&comment.ID)
	if err != nil {
		return nil, fmt.Errorf("inserting comment: %w", err)
	}
	return comment, nil
}

func (l *DBLoader) DeleteComment(commentID int) error {
	_, err := l.db.Exec("DELETE FROM comments WHERE id = $1", commentID)
	if err != nil {
		return fmt.Errorf("deleting comment: %w", err)
	}
	return nil
}
