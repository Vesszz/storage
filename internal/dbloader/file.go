package dbloader

import (
	"fmt"
	"go.uber.org/zap"
	"main/internal/models"
)

func (l *DBLoader) InsertFile(fileModel *models.File) (*models.File, error) {
	err := l.db.QueryRow("INSERT INTO files(user_id, key, path, time_created, name, description, times_viewed) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id",
		fileModel.UserID, fileModel.Key, fileModel.Path, fileModel.TimeCreated, fileModel.Name, fileModel.Description, fileModel.TimesViewed).Scan(&fileModel.ID)
	if err != nil {
		zap.S().Errorf("insert file: %v", err)
		return nil, fmt.Errorf("insert file: %w", err)
	}
	return fileModel, err
}

func (l *DBLoader) DeleteFile(fileID int) error {
	_, err := l.db.Exec("DELETE FROM files WHERE id = $1", fileID)
	if err != nil {
		zap.S().Errorf("delete file: %v", err)
		return fmt.Errorf("deleting file: %w", err)
	}
	return nil
}
