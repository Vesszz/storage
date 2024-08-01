package loader

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"main/internal/config"
)

type Loader struct {
	dbCfg *config.DatabaseConfig
	fsCfg *config.FileStorageConfig
	db    *sql.DB
}

func New(dbCfg *config.DatabaseConfig, fsCfg *config.FileStorageConfig) (*Loader, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s", dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.Host, dbCfg.Port, dbCfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening db connection: %w", err)
	}
	return &Loader{
		dbCfg: dbCfg,
		fsCfg: fsCfg,
		db:    db,
	}, nil
}
