package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBUrl          string
	AuthServiceURL string
	UserServiceURL string
	JWTSecret      string
}

func LoadConfig() *Config {
	_ = godotenv.Load(".env")

	return &Config{
		Port:           getEnv("PORT", "8000"),
		DBUrl:          getEnv("DB_URL", ""),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", ""),
		UserServiceURL: getEnv("USER_SERVICE_URL", ""),
		JWTSecret:      getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
