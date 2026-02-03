package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"name"`
	Port     int    `yaml:"port"`
	SSLMode  string `yaml:"sslmode"`
}

func LoadConfig() (*Config, error) {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
