package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App AppConfig
}

type AppConfig struct {
	Port            int
	Environment     string
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int
}

type Environment string

const (
	Development Environment = "development"
	Stage       Environment = "stage"
	Production  Environment = "production"
)

const (
	EnvKeyPort             = "APP_PORT"
	EnvKeyEnvironment      = "APP_ENVIRONMENT"
	EnvKeyReadTimeout      = "APP_READ_TIMEOUT"
	EnvKeyWriteTimeout     = "APP_WRITE_TIMEOUT"
	EenvKeyShutdownTimeout = "APP_SHUTDOWN_TIMEOUT"
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

	shutdownTimeout, err := getEnvAsInt(EenvKeyShutdownTimeout, DefaultShutdownTimeout)
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

func LoadConfig() (*Config, error) {
	if err := loadEnvFile(); err != nil {
		return nil, fmt.Errorf("could not load .env file: %v", err)
	}

	appConfig, err := loadAppConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load app config: %w", err)
	}

	return &Config{
		App: appConfig,
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
