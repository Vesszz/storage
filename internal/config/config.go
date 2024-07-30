package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Database    DatabaseConfig    `yaml:"database"`
	FileStorage FileStorageConfig `yaml:"file_storage"`
	Server      ServerConfig      `yaml:"server"`
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
