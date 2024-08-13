package logic

import (
	"fmt"
	"main/internal/errs"
	"main/internal/models"
	"main/pkg/secutils"
	"time"
)

func (l *Logic) Register(name string, password string, fingerprint string) (*UserTokens, error) {
	user, err := l.dbLoader.GetUserByName(name)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}
	if user != nil {
		return nil, errs.UserAlreadyExists
	}

	hashedPassword, err := secutils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hashing password: %w", err)
	}
	user, err = l.dbLoader.CreateUser(models.User{
		Name:        name,
		Password:    hashedPassword,
		TimeCreated: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("registration on db side: %w", err)
	}

	userTokens, err := l.CreateTokens(user, fingerprint)
	if err != nil {
		return nil, fmt.Errorf("creating tokens: %w", err)
	}
	return userTokens, nil
}

func (l *Logic) Login(name string, password string, fingerprint string) (*UserTokens, error) {
	user, err := l.dbLoader.GetUserByName(name)
	if err != nil {
		return nil, fmt.Errorf("finding user: %w", err)
	}
	if user == nil {
		return nil, errs.UserNotExists
	}
	hashedPassword, err := secutils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hashing while signing in: %w", err)
	}
	if !secutils.CompareHashAndPassword(hashedPassword, password) {
		return nil, errs.WrongPassword
	}
	userTokens, err := l.CreateTokens(user, fingerprint)
	if err != nil {
		return nil, fmt.Errorf("creating tokens: %w", err)
	}
	return userTokens, nil
}

func (l *Logic) Refresh(refreshToken string, fingerprint string) (*UserTokens, error) {
	rtm, err := l.dbLoader.GetRefreshTokenModelByRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("get refresh token: %w", err)
	}
	err = l.CheckRefreshToken(*rtm, fingerprint)
	if err != nil {
		return nil, fmt.Errorf("check refresh token: %w", err)
	}
	user, err := l.dbLoader.GetUserByID(rtm.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	tokens, err := l.CreateTokens(user, fingerprint)
	if err != nil {
		return nil, fmt.Errorf("create tokens: %w", err)
	}
	return tokens, nil
}
