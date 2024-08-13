package logic

import (
	"fmt"
	"main/internal/config"
	"main/internal/dbloader"
	"main/internal/fileloader"
	"os"
)

type Logic struct {
	fileLoader     *fileloader.FileLoader
	dbLoader       *dbloader.DBLoader
	config         *config.Config
	secretTokenKey string
}

func New(c *config.Config) (*Logic, error) {
	secretTokenKey, ok := os.LookupEnv("SECRET_TOKEN_KEY")
	if !ok {
		return nil, fmt.Errorf("loading env")
	}
	fl, err := fileloader.New(&c.FileStorage)
	if err != nil {
		return nil, fmt.Errorf("creating fileloader: %w", err)
	}
	dbl, err := dbloader.New(&c.Database)
	if err != nil {
		return nil, fmt.Errorf("creating dbloader: %w", err)
	}
	return &Logic{
		fileLoader:     fl,
		dbLoader:       dbl,
		config:         c,
		secretTokenKey: secretTokenKey,
	}, nil
}
