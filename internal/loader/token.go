package loader

import (
	"database/sql"
	"errors"
	"fmt"
	"main/internal/errs"
	"main/internal/models"
)

func (l *Loader) UpdateRefreshToken(refreshToken models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens(user_id, fingerprint, key, expires_at) 
		VALUES ($1, $2, $3, $4) 
		ON CONFLICT (user_id, fingerprint) 
		DO UPDATE SET 
    		key = EXCLUDED.key, 
    		expires_at = EXCLUDED.expires_at
	`
	_, err := l.db.Exec(query, refreshToken.UserID, refreshToken.Fingerprint, refreshToken.Key, refreshToken.ExpiresAt)
	if err != nil {
		return fmt.Errorf("inserting refresh token: %w", err)
	}
	return nil
}

func (l *Loader) GetRefreshToken(userID int64) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := l.db.QueryRow("SELECT user_id, fingerprint, key, expires_at FROM refresh_tokens WHERE user_id=$1", userID).Scan(&refreshToken.UserID, &refreshToken.Fingerprint, &refreshToken.Key, &refreshToken.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.UserNotExists
		} else {
			return nil, fmt.Errorf("selecting user: %w", err)
		}
	}
	return &refreshToken, nil
}

func (l *Loader) GetRefreshTokenByRefreshToken(refreshToken string) (*models.RefreshToken, error) {
	var refreshTokenModel models.RefreshToken
	err := l.db.QueryRow("SELECT user_id, fingerprint, key, expires_at FROM refresh_tokens WHERE key=$1", refreshToken).
		Scan(&refreshTokenModel.UserID, &refreshTokenModel.Fingerprint, &refreshTokenModel.Key, &refreshTokenModel.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errs.UserNotExists
		} else {
			return nil, fmt.Errorf("selecting user: %w", err)
		}
	}
	return &refreshTokenModel, nil
}
