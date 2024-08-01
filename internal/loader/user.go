package loader

import (
	"database/sql"
	"errors"
	"fmt"
	"main/internal/models"
)

func (l *Loader) CreateUser(user models.User) (*models.User, error) {
	err := l.db.QueryRow("INSERT INTO users(name, password, salt) VALUES ($1, $2, $3) RETURNING id", user.Name, user.Password, user.Salt).Scan(&user.ID)
	if err != nil {
		return nil, fmt.Errorf("inserting user: %w", err)
	}
	return &user, nil
}

func (l *Loader) FindUser(userName string) (*models.User, error) {
	var user models.User
	err := l.db.QueryRow("SELECT id, name, password, salt FROM users WHERE name=$1", userName).Scan(&user.ID, &user.Name, &user.Password, &user.Salt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, fmt.Errorf("selecting user: %w", err)
		}
	}
	return &user, nil
}
