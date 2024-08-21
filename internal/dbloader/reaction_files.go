package dbloader

import (
	"fmt"
	"go.uber.org/zap"
	"main/internal/models"
)

func (l *DBLoader) InsertFileReaction(fileReaction models.ReactionFile) error {
	_, err := l.db.Exec("INSERT INTO reaction_files(user_id, file_id, type_id) VALUES ($1, $2, $3)",
		fileReaction.UserID, fileReaction.FileID, fileReaction.TypeID)
	if err != nil {
		zap.S().Errorf("insert file reaction: %v", err)
		return fmt.Errorf("insert file reaction: %w", err)
	}
	return nil
}

func (l *DBLoader) DeleteFileReaction(fileReaction models.ReactionFile) error {
	_, err := l.db.Exec("DELETE FROM reaction_files WHERE user_id = $1 AND file_id = $2",
		fileReaction.UserID, fileReaction.FileID)
	if err != nil {
		zap.S().Errorf("delete file reaction: %v", err)
		return fmt.Errorf("delete file reaction: %w", err)
	}
	return nil
}
