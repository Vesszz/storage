package logic

import (
	"fmt"
	"io"
	"log/slog"
	"main/internal/config"
	"main/internal/errs"
	"main/internal/loader"
	"main/internal/models"
	"main/pkg/secutils"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type Logic struct {
	loader         *loader.Loader
	config         *config.Config
	secretTokenKey string
}

func New(c *config.Config) (*Logic, error) {
	secretTokenKey, ok := os.LookupEnv("SECRET_TOKEN_KEY")
	if !ok {
		return nil, fmt.Errorf("loading env")
	}
	l, err := loader.New(&c.Database, &c.FileStorage)
	if err != nil {
		return nil, fmt.Errorf("loader initialisation: %w", err)
	}
	return &Logic{
		loader:         l,
		config:         c,
		secretTokenKey: secretTokenKey,
	}, nil
}

func (l *Logic) Index() (loader.PageData, error) {
	indexInfo, err := l.loader.IndexInfo()
	if err != nil {
		return loader.PageData{}, fmt.Errorf("getting page info: %w", err)
	}
	return indexInfo, nil
}

func (l *Logic) Upload(fileName string, file multipart.File, atc *AccessTokenClaims) error {
	_, err := l.loader.SaveFile(models.File{
		UserID:      atc.ID,
		TimeCreated: time.Now(),
		Name:        fileName,
		TimesViewed: 0,
	}, file)
	if err != nil {
		slog.Error("uploading file: %w", err)
		return fmt.Errorf("uploading file: %w", err)
	}
	return nil
}

func (l *Logic) Download(w http.ResponseWriter, fileName string) error {
	file, err := l.loader.Get("./uploads/" + fileName)
	if err != nil {
		return fmt.Errorf("getting file: %w", err)
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, file)
	if err != nil {
		return fmt.Errorf("downloading: %w", err)
	}
	defer file.Close()
	return nil
}

func (l *Logic) FileList() (*loader.FileList, error) {
	files, err := l.loader.GetAll()
	if err != nil {
		return nil, fmt.Errorf("getting all files %w", err)
	}
	return files, nil
}

func (l *Logic) Register(name string, password string, fingerprint string) (*UserTokens, error) {
	user, err := l.loader.GetUserByName(name)
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
	user, err = l.loader.CreateUser(models.User{
		Name:     name,
		Password: hashedPassword,
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
	user, err := l.loader.GetUserByName(name)
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
	rtm, err := l.loader.GetRefreshTokenByRefreshToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("get refresh token: %w", err)
	}
	err = l.CheckRefreshToken(*rtm, fingerprint)
	if err != nil {
		return nil, fmt.Errorf("check refresh token: %w", err)
	}
	user, err := l.loader.GetUserByID(rtm.UserID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	tokens, err := l.CreateTokens(user, fingerprint)
	if err != nil {
		return nil, fmt.Errorf("create tokens: %w", err)
	}
	return tokens, nil
}
