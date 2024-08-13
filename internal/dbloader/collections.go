package dbloader

import (
	"fmt"
	"main/internal/models"
)

func (l *DBLoader) CreateCollection(collection *models.Collection) (*models.Collection, error) {
	err := l.db.QueryRow("INSERT INTO collections(user_id, name) VALUES ($1, $2) RETURNING id", collection.UserID, collection.Name).Scan(&collection.ID)
	if err != nil {
		return nil, fmt.Errorf("inserting collection: %w", err)
	}
	return collection, nil
}

func (l *DBLoader) DeleteCollection(collectionID int) error {
	_, err := l.db.Exec("DELETE FROM collections WHERE id = $1", collectionID)
	if err != nil {
		return fmt.Errorf("deleting file from collection: %w", err)
	}
	return nil
}

func (l *DBLoader) AddFileToCollection(collectionFile *models.CollectionFile) error {
	_, err := l.db.Exec("INSERT INTO collection_files(collection_id, file_id) VALUES ($1, $2)", collectionFile.CollectionID, collectionFile.FileID)
	if err != nil {
		return fmt.Errorf("inserting file into collection: %w", err)
	}
	return nil
}

func (l *DBLoader) RemoveFileFromCollection(collectionFile *models.CollectionFile) error {
	_, err := l.db.Exec("DELETE FROM collection_files WHERE collection_id = $1 AND file_id = $2", collectionFile.CollectionID, collectionFile.FileID)
	if err != nil {
		return fmt.Errorf("deleting file from collection: %w", err)
	}
	return nil
}
