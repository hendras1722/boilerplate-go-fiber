package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv    string
	AppPort   string
	DBHost    string
	DBPort    string
	DBUser    string
	DBPass    string
	DBName    string
	JWTSecret         string
	JWTRefreshSecret  string
	JWTExpHours       string
	JWTRefreshExpHours string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	return &Config{
		AppEnv:    getEnv("APP_ENV", "development"),
		AppPort:   getEnv("APP_PORT", "3000"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBPort:    getEnv("DB_PORT", "5432"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASSWORD", "password"),
		DBName:             getEnv("DB_NAME", "app_db"),
		JWTSecret:          getEnv("JWT_SECRET", "super-secret-key"),
		JWTRefreshSecret:   getEnv("JWT_REFRESH_SECRET", "super-refresh-secret-key"),
		JWTExpHours:        getEnv("JWT_EXP_HOURS", "24"),
		JWTRefreshExpHours: getEnv("JWT_REFRESH_EXP_HOURS", "168"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
