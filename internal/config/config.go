package config

import (
	"fmt"
	"os"
)

type Config struct {
	Environment string
	DBHost      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBPort      string
	JwtSecret   string
}

func NewConfig() *Config {
	return &Config{
		Environment: getEnvOrDefault("ENVIRONMENT", "develop"),
		DBHost:      getEnvOrDefault("DB_HOST", "localhost"),
		DBUser:      getEnvOrDefault("DB_USER", "postgres"),
		DBPassword:  getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:      getEnvOrDefault("DB_NAME", "go_crud"),
		DBPort:      getEnvOrDefault("DB_PORT", "5432"),
		JwtSecret:   getEnvOrDefault("JWT_SECRET", ""),
	}
}

// getEnvOrDefault returns the value of an environment variable or a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvironment() string {
	config := NewConfig()
	return config.Environment
}

func GetDBConnectionString() string {
	config := NewConfig()
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
}

func GetJwtSecret() []byte {
	config := NewConfig()
	if config.JwtSecret == "" {
		panic("JWT_SECRET environment variable is not set. Please set a strong JWT secret.")
	}
	if len(config.JwtSecret) < 32 {
		panic("JWT_SECRET must be at least 32 characters long for security.")
	}
	return []byte(config.JwtSecret)
}
