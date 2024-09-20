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

	// Какая? Redis или memcached
	Database string
	// Номер БД (юзаем один образ для двух сервисов, надо выбрать разные бд)
	Number string
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
			Host:     os.Getenv("MEMORY_DB_HOST"),
			Port:     os.Getenv("MEMORY_DB_PORT"),
			Database: os.Getenv("MEMORY_DB"),
			Number:   os.Getenv("MEMORY_DB_NUMBER"),
		},
	}, nil
}
