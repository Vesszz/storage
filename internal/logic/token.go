package logic

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"main/internal/errs"
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

type AccessTokenClaims struct {
	Name      string `json:"name"`
	ExpiresAt int64  `json:"expires_at"`
	ID        int64  `json:"id"`
}

func (atc *AccessTokenClaims) ToMapStringInterface() (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(atc)
	if err != nil {
		return nil, fmt.Errorf("marshaling access token claims: %w", err)
	}
	var msi map[string]interface{}
	err = json.Unmarshal(jsonBytes, &msi)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling access token claims into map: %w", err)
	}
	return msi, nil
}

func (l *Logic) CreateTokens(user *models.User, fingerprint string) (*UserTokens, error) {
	accessToken, err := l.CreateAccessToken(user)
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

func (l *Logic) CreateAccessToken(user *models.User) (string, error) {
	//todo access token fields
	atc := AccessTokenClaims{
		Name:      user.Name,
		ExpiresAt: time.Now().Add(AccessTokenLifetime).Unix(),
		ID:        user.ID,
	}
	atcmap, err := atc.ToMapStringInterface()
	if err != nil {
		return "", fmt.Errorf("changing access token claims into map: %w", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(atcmap))
	tokenString, err := token.SignedString([]byte(l.secretTokenKey))
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}
	return tokenString, nil
}

func (l *Logic) ValidAccessToken(accessToken string, user *models.User) (*AccessTokenClaims, error) {
	atc, err := l.AccessTokenClaimsFromAccessToken(accessToken)
	if err != nil {
		return nil, errs.InvalidAccessToken
	}
	if atc.ExpiresAt < time.Now().Unix() {
		return nil, errs.AccessTokenExpired
	}
	if atc.Name != user.Name {
		return nil, errs.InvalidAccessToken
	}
	if atc.ID != user.ID {
		return nil, errs.InvalidAccessToken
	}

	return atc, nil
}

func (l *Logic) CheckRefreshToken(rtm models.RefreshToken, fingerprint string) error {
	if rtm.Fingerprint != fingerprint {
		return errs.InvalidRefreshToken
	}
	if rtm.ExpiresAt.Unix() < time.Now().Unix() {
		return errs.RefreshTokenExpired
	}
	return nil
}

func (l *Logic) AccessTokenClaimsFromAccessToken(accessToken string) (*AccessTokenClaims, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		//if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		//	return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		//}
		return []byte(l.secretTokenKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errs.InvalidAccessToken
	}
	atc, err := AccessTokenClaimsFromJWTMapClaims(claims)
	if err != nil {
		return nil, fmt.Errorf("access token claims from jwt map claims: %w", err)
	}
	return atc, nil
}

func AccessTokenClaimsFromJWTMapClaims(claims jwt.MapClaims) (*AccessTokenClaims, error) {
	var act AccessTokenClaims
	jsonBytes, err := json.Marshal(claims)
	if err != nil {
		return nil, fmt.Errorf("marshal jwt map claims: %w", err)
	}
	err = json.Unmarshal(jsonBytes, &act)
	if err != nil {
		return nil, fmt.Errorf("unmarshal jwt map claims to access token claims: %w", err)
	}
	return &act, nil
}
