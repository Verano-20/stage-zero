package config

import (
	"fmt"
	"os"
	"time"
)

var globalConfig *Config

type Config struct {
	ServiceName    string
	ServiceVersion string
	ServicePort    string
	Environment    string
	JwtSecret      string
	Database       DatabaseConfig
	Telemetry      TelemetryConfig
}

type DatabaseConfig struct {
	Host     string
	User     string
	Password string
	Name     string
	Port     string
}

type TelemetryConfig struct {
	EnableStdout   bool
	EnableOTLP     bool
	OTLPEndpoint   string
	OTLPInsecure   bool
	MetricInterval time.Duration
}

func InitConfig() {
	config := &Config{
		ServiceName:    getEnvOrDefault("SERVICE_NAME", "stage-zero-api"),
		ServiceVersion: getEnvOrDefault("SERVICE_VERSION", "1.0.0"),
		ServicePort:    getEnvOrDefault("SERVICE_PORT", "8080"),
		Environment:    getEnvOrDefault("ENVIRONMENT", "develop"),
		JwtSecret:      getEnvOrDefault("JWT_SECRET", ""),
		Database:       *initDatabaseConfig(),
		Telemetry:      *initTelemetryConfig(),
	}

	globalConfig = config
}

func initDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
		Name:     getEnvOrDefault("DB_NAME", "go_crud"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
	}
}

func initTelemetryConfig() *TelemetryConfig {
	metricInterval, err := time.ParseDuration(getEnvOrDefault("METRIC_INTERVAL", "30s"))
	if err != nil {
		panic("Invalid METRIC_INTERVAL: " + err.Error())
	}

	return &TelemetryConfig{
		EnableStdout:   getEnvOrDefault("ENABLE_STDOUT", "true") == "true",
		EnableOTLP:     getEnvOrDefault("ENABLE_OTLP", "true") == "true",
		OTLPEndpoint:   getEnvOrDefault("OTLP_ENDPOINT", "localhost:4317"),
		OTLPInsecure:   getEnvOrDefault("OTLP_INSECURE", "true") == "true",
		MetricInterval: metricInterval,
	}
}

func Get() *Config {
	if globalConfig == nil {
		panic("Config not initialized")
	}
	return globalConfig
}

func (config *Config) GetDBConnectionString() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.Database.Host, config.Database.User, config.Database.Password, config.Database.Name, config.Database.Port)
}

func (config *Config) GetJwtSecret() []byte {
	if config.JwtSecret == "" {
		panic("JWT_SECRET environment variable is not set. Please set a strong JWT secret.")
	}
	if len(config.JwtSecret) < 32 {
		panic("JWT_SECRET must be at least 32 characters long for security.")
	}
	return []byte(config.JwtSecret)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
