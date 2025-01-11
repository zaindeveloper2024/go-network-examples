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
	Port        int
	Environment string
}

const (
	DefaultAppPort        = 8080
	DefaultAppEnvironment = "development"
)

func LoadConfig() (*Config, error) {
	godotenv.Load()
	fmt.Println("Loading .env file")

	config := Config{}

	port, err := getEnvAsInt("APP_PORT", DefaultAppPort)
	if err != nil {
		return nil, fmt.Errorf("could not get PORT: %v", err)
	}

	config.App = AppConfig{
		Port:        port,
		Environment: getEnv("APP_ENVIRONMENT", DefaultAppEnvironment),
	}

	return &config, nil
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
