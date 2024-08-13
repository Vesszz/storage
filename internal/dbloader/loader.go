package dbloader

import (
	"database/sql"
	"fmt"
	"main/internal/config"

	_ "github.com/lib/pq"
)

type DBLoader struct {
	dbCfg *config.DatabaseConfig
	db    *sql.DB
}

func New(dbCfg *config.DatabaseConfig) (*DBLoader, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s", dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.Host, dbCfg.Port, dbCfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening db connection: %w", err)
	}
	return &DBLoader{
		dbCfg: dbCfg,
		db:    db,
	}, nil
}
