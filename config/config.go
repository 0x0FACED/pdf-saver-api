package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PDF      PRFServiceConfig
	MemCache MemCacheConfig
}

type MemCacheConfig struct {
	Host string
	Port string
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
		MemCache: MemCacheConfig{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
		},
	}, nil
}
