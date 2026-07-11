package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	LogLevel    string
	AppPort     string
	JWTSecret   string
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
}

func Load() Config {
	_ = godotenv.Load()
	return Config{
		Environment: valueOrDefault("API_ENV", "development"),
		LogLevel:    valueOrDefault("LOG_LEVEL", "info"),
		AppPort:     valueOrDefault("APP_PORT", "8080"),
		JWTSecret:   valueOrDefault("JWT_SECRET", "development-secret"),
		DBHost:      valueOrDefault("DB_HOST", "127.0.0.1"),
		DBPort:      valueOrDefault("DB_PORT", "5432"),
		DBUser:      valueOrDefault("DB_USER", "postgres"),
		DBPassword:  valueOrDefault("DB_PASS", "password"),
		DBName:      valueOrDefault("DB_NAME", "withdraw_db"),
		DBSSLMode:   valueOrDefault("DB_SSL_MODE", "disable"),
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode)
}

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
