package config

import (
	"log"
	"os"
)

type Config struct {
	DBUrl      string
	JWTSecret  string
	ServerPort string
}

func LoadConfig() *Config {
	cfg := &Config{
		DBUrl:      os.Getenv("DATABASE_URL"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		ServerPort: os.Getenv("PORT"),
	}

	if cfg.DBUrl == "" || cfg.JWTSecret == "" {
		log.Fatalf("Missing environment variables")
	}
	return cfg
}
