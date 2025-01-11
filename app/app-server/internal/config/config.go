package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App AppConfig
	DB  DBConfig
}

type AppConfig struct {
	Port            int
	Environment     Environment
	ReadTimeout     int
	WriteTimeout    int
	ShutdownTimeout int
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
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
	EnvKeyDBHost          = "DB_HOST"
	EnvKeyDBPort          = "DB_PORT"
	EnvKeyDBUser          = "DB_USER"
	EnvKeyDBPassword      = "DB_PASSWORD"
	EnvKeyDBName          = "DB_NAME"
	EnvKeyDBSSLMode       = "DB_SSL_MODE"
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

	env, err := parseEnvironment(os.Getenv(EnvKeyEnvironment))
	if err != nil {
		return AppConfig{}, fmt.Errorf("could not get APP_ENVIRONMENT: %v", err)
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
		Environment:     env,
		ReadTimeout:     readTimeout,
		WriteTimeout:    writeTimeout,
		ShutdownTimeout: shutdownTimeout,
	}, nil
}

func loadDBConfig() (DBConfig, error) {
	port, err := getEnvAsInt(EnvKeyDBPort, 5432)
	if err != nil {
		return DBConfig{}, fmt.Errorf("could not get DB_PORT: %v", err)
	}

	return DBConfig{
		Host:     getEnv(EnvKeyDBHost, "localhost"),
		Port:     port,
		User:     getEnv(EnvKeyDBUser, "postgres"),
		Password: getEnv(EnvKeyDBPassword, "password"),
		Name:     getEnv(EnvKeyDBName, "dbname"),
		SSLMode:  getEnv(EnvKeyDBSSLMode, "disable"),
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
		App: appConfig,
		DB:  databaseConfig,
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

func parseEnvironment(env string) (Environment, error) {
	switch Environment(env) {
	case Development, Stage, Production:
		return Environment(env), nil
	default:
		return "", fmt.Errorf("invalid environment: %s", env)
	}
}

func (c *Config) IsDevelopment() bool {
	return c.App.Environment == Development
}

func (c *DBConfig) DNS() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode)
}
