package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

var Cfg Config

type Config struct {
	App          AppConfig
	Database     DatabaseConfig
	DatabaseTest DatabaseConfig
	Api          ApiConfig
}

type AppConfig struct {
	Port         int
	JwtSecret    string
	LogFile      string
	Environment  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type DatabaseConfig struct {
	Host         string
	Port         string
	Name         string
	User         string
	Password     string
	SSLMode      string
	MigrationDir string
}

type ApiConfig struct {
	Limiter   *rate.Limiter
	UserAgent string
}

func MustLoad(path string) error {
	if path == "" {
		path = ".env"
	}
	err := godotenv.Load(path)
	if err != nil {
		return err
	}
	portStr := getEnv("APP_PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return err
	}

	Cfg = Config{
		App: AppConfig{
			Port: port,
			// 6. Безопасность
			// - JWT secret хранить в env
			JwtSecret: getEnv("JWT_SECRET", "dev-secret-change-me"),
			// Ch 2. Logging Lv 5. Logger Configuration
			// Assume that in production,
			// WeatherApp has a LOG_FILE environment variable set.
			// In local development and staging, it is not set.
			LogFile:      getEnv("LOG_FILE", ""),
			Environment:  getEnv("ENV", "development"),
			ReadTimeout:  mustDuration("APP_READ_TIMEOUT", "5s"),
			WriteTimeout: mustDuration("APP_WRITE_TIMEOUT", "10s"),
			IdleTimeout:  mustDuration("APP_IDLE_TIMEOUT", "60s"),
		},
		Database: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			Name:         getEnv("DB_NAME", "users_db"),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", "postgres"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MigrationDir: getEnv("MIGRATION_DIR", "./migrations/postgres"),
		},
		DatabaseTest: DatabaseConfig{
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "5432"),
			Name:         getEnv("DB_NAME_TEST", "users_db_test"),
			User:         getEnv("DB_USER", "postgres"),
			Password:     getEnv("DB_PASSWORD", "postgres"),
			SSLMode:      getEnv("DB_SSLMODE", "disable"),
			MigrationDir: getEnv("MIGRATION_DIR", "./migrations/postgres"),
		},
		Api: ApiConfig{
			Limiter:   rate.NewLimiter(rate.Every(time.Second), 5),
			UserAgent: getEnv("USER_AGENT", "weather-api/1.0 (example@gmail.com)"),
		},
	}
	return nil
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Name,
		c.User,
		c.Password,
		c.SSLMode,
	)
}

////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions
////// accommodating functions

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func mustDuration(key, fallback string) time.Duration {
	value := getEnv(key, fallback)
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("invalid duration %s=%s: %v", key, value, err)
	}
	return d
}
