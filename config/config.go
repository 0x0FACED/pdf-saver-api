package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PDF PRFServiceConfig
}

type PRFServiceConfig struct {
	Host string
	Port string
}

func Load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		PDF: PRFServiceConfig{
			Host: os.Getenv("PDF_HOST"),
			Port: os.Getenv("PDF_PORT"),
		},
	}, nil
}
