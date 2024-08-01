package logic

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"main/internal/models"
	"time"
)

const (
	AccessTokenLifetime  = time.Minute * 30
	RefreshTokenLifetime = time.Hour * 24 * 30
)

type UserTokens struct {
	RefreshToken uuid.UUID
	AccessToken  string
}

func (l *Logic) CreateTokens(user *models.User, fingerprint string) (*UserTokens, error) {
	accessToken, err := l.createAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("creating access token: %w", err)
	}

	refreshToken := uuid.New()
	err = l.loader.UpdateRefreshToken(models.RefreshToken{
		Key:         refreshToken,
		UserID:      user.ID,
		Fingerprint: fingerprint,
		ExpiresAt:   time.Now().Add(RefreshTokenLifetime),
	})
	if err != nil {
		return nil, fmt.Errorf("updating refresh token: %w", err)
	}

	return &UserTokens{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

func (l *Logic) createAccessToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Name,
			"exp":      time.Now().Add(AccessTokenLifetime).Unix(),
		})
	tokenString, err := token.SignedString([]byte(l.secretTokenKey))
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}
	return tokenString, nil
}
