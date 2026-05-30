package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	HTTP  HTTP
	Redis Redis
}

type HTTP struct {
	Host string
	Port string
}

type Redis struct {
	Addr     string
	Password string
	DB       int
}

func LoadConfig() (*AppConfig, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	config := &AppConfig{
		HTTP: HTTP{
			Host: os.Getenv("HTTP_HOST"),
			Port: os.Getenv("HTTP_PORT"),
		},

		Redis: Redis{
			Addr:     os.Getenv("REDIS_HOST"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       db,
		},
	}

	return config, nil
}
