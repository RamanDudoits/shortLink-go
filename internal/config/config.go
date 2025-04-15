package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DB struct {
		DSN string
	}
	HTTP struct {
		Addr string
		ReadTimeout time.Duration
		WriteTimeout time.Duration
	}
	JWT struct {
		Secret string
		Expire time.Duration
	}
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	var cfg Config

	cfg.DB.DSN = os.Getenv("DB_DSN")
	if cfg.DB.DSN == "" {
		cfg.DB.DSN = "postgres://user:password@localhost:5432/shortLink-go?sslmode=disable"
		
	}
	cfg.HTTP.Addr = os.Getenv("HTTP_ADDR")
	if cfg.HTTP.Addr == "" {
		cfg.HTTP.Addr = ":8080"
	}
	cfg.HTTP.ReadTimeout = 10 * time.Second
	cfg.HTTP.WriteTimeout = 10 * time.Second

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	if cfg.JWT.Secret == "" {
		cfg.JWT.Secret = "GYHUkgeBkqfC"
	}
	cfg.JWT.Expire = 24 * time.Hour

	return &cfg, nil
}