package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Port            int
	Environment     string
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int
}

type DatabaseConfig struct {
	URL string
}

type Environment string

const (
	Development Environment = "development"
	Stage       Environment = "stage"
	Production  Environment = "production"
)

const (
	EnvKeyPort            = "APP_PORT"
	EnvKeyEnvironment     = "APP_ENVIRONMENT"
	EnvKeyReadTimeout     = "APP_READ_TIMEOUT"
	EnvKeyWriteTimeout    = "APP_WRITE_TIMEOUT"
	EnvKeyShutdownTimeout = "APP_SHUTDOWN_TIMEOUT"
	EnvKeyDatabaseURL     = "DATABASE_URL"
)

const (
	DefaultAppPort         = 8080
	DefaultAppEnvironment  = "development"
	DefaultReadTimeout     = 30
	DefaultWriteTimeout    = 30
	DefaultShutdownTimeout = 30
)

func loadEnvFile() error {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("could not load .env file: %v", err)
		}
	}
	return nil
}

func loadAppConfig() (AppConfig, error) {
	port, err := getEnvAsInt(EnvKeyPort, DefaultAppPort)
	if err != nil {
		return AppConfig{}, fmt.Errorf("could not get APP_PORT: %v", err)
	}

	readTimeout, err := getEnvAsInt(EnvKeyReadTimeout, DefaultReadTimeout)
	if err != nil {
		return AppConfig{}, fmt.Errorf("could not get APP_READ_TIMEOUT: %v", err)
	}

	writeTimeout, err := getEnvAsInt(EnvKeyWriteTimeout, DefaultWriteTimeout)
	if err != nil {
		return AppConfig{}, fmt.Errorf("could not get APP_WRITE_TIMEOUT: %v", err)
	}

	shutdownTimeout, err := getEnvAsInt(EnvKeyShutdownTimeout, DefaultShutdownTimeout)
	if err != nil {
		return AppConfig{}, fmt.Errorf("could not get APP_SHUTDOWN_TIMEOUT: %v", err)
	}

	return AppConfig{
		Port:            port,
		Environment:     getEnv(EnvKeyEnvironment, DefaultAppEnvironment),
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		ShutdownTimeout: shutdownTimeout,
	}, nil
}

func loadDBConfig() (DatabaseConfig, error) {
	databaseUrl := os.Getenv("DATABASE_URL")
	if databaseUrl == "" {
		return DatabaseConfig{}, fmt.Errorf("DATABASE_URL is required")
	}

	return DatabaseConfig{
		URL: databaseUrl,
	}, nil
}

func LoadConfig() (*Config, error) {
	if err := loadEnvFile(); err != nil {
		return nil, fmt.Errorf("could not load .env file: %v", err)
	}

	appConfig, err := loadAppConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load app config: %w", err)
	}

	databaseConfig, err := loadDBConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load database config: %w", err)
	}

	return &Config{
		App:      appConfig,
		Database: databaseConfig,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) (int, error) {
	if value, exists := os.LookupEnv(key); exists {
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return 0, fmt.Errorf("invalid value for %s: %s", key, err)
		}
		return intVal, nil
	}
	return defaultValue, nil
}
