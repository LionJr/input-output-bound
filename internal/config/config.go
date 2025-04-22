package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
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
	Db       int
}

func LoadConfig() (*AppConfig, error) {
	err := godotenv.Load()
	if err != nil {
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
			Db:       db,
		},
	}

	return config, nil
}
