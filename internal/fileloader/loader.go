package fileloader

import "main/internal/config"

type FileLoader struct {
	fsCfg *config.FileStorageConfig
}

func New(fsCfg *config.FileStorageConfig) (*FileLoader, error) {
	return &FileLoader{
		fsCfg: fsCfg,
	}, nil
}
