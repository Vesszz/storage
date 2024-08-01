package loader

import (
	"fmt"
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
