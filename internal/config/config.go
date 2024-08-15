package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Database    DatabaseConfig    `yaml:"database"`
	FileStorage FileStorageConfig `yaml:"file_storage"`
	Server      ServerConfig      `yaml:"server"`
	Logger      LoggerConfig      `yaml:"logger"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type FileStorageConfig struct {
	Path string `yaml:"path"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type LoggerConfig struct {
	Level string            `yaml:"level"`
	File  *FileLoggerConfig `yaml:"file"`
}

type FileLoggerConfig struct {
	Filename   string        `yaml:"filename"`
	MaxSize    int64         `yaml:"max_size"`
	MaxAge     time.Duration `yaml:"max_age"`
	MaxBackups int           `yaml:"max_backups"`
	LocalTime  bool          `yaml:"localtime"`
	Compress   bool          `yaml:"compress"`
}

func Load(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}
	var config Config
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshaling config file: %w", err)
	}
	return &config, nil
}
