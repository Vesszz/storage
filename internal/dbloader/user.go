package dbloader

import (
	"database/sql"
	"errors"
	"fmt"
	"main/internal/models"
)

func (l *DBLoader) CreateUser(user models.User) (*models.User, error) {
	err := l.db.QueryRow("INSERT INTO users(name, password, time_created) VALUES ($1, $2, $3) RETURNING id", user.Name, user.Password, user.TimeCreated).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("inserting user: %w", err)
	}
	return &user, nil
}

func (l *DBLoader) GetUserByName(userName string) (*models.User, error) {
	var user models.User
	err := l.db.QueryRow("SELECT id, name, password FROM users WHERE name=$1", userName).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("selecting user: %w", err)
		}
	}
	return &user, nil
}

func (l *DBLoader) GetUserByID(userID int64) (*models.User, error) {
	var user models.User
	err := l.db.QueryRow("SELECT id, name, password FROM users WHERE id=$1", userID).Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("selecting user: %w", err)
		}
	}
	return &user, nil
}
