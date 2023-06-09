package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Auth     AuthConfig
}

type AppConfig struct {
	Env  string
	Port int
}

type AuthConfig struct {
	ExpireDuration time.Duration
	SecretKey      string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

func NewConfig() *Config {
	return &Config{
		App: AppConfig{
			Env:  GetEnv("APP_ENV", "development"),
			Port: GetEnvInt("APP_PORT", 9999),
		},
		Auth: AuthConfig{
			ExpireDuration: GetTimeDuration("AUTH_EXPIRE_DURATION", 1*time.Hour),
			SecretKey:      GetEnv("AUTH_SECRET_KEY", "secret"),
		},
		Database: DatabaseConfig{
			Host:     GetEnv("DATABASE_HOST", "localhost"),
			Port:     GetEnvInt("DATABASE_PORT", 5432),
			User:     GetEnv("DATABASE_USER", "postgres"),
			Password: GetEnv("DATABASE_PASSWORD", "postgres"),
			Name:     GetEnv("DATABASE_NAME", "todo"),
		},
	}
}

func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvInt(key string, fallback int) int {
	if str, ok := os.LookupEnv(key); ok {
		value, err := strconv.Atoi(str)
		if err != nil {
			log.Fatalf("bad value converting %s to int: %v", key, err)
		}
		return value
	}
	return fallback
}

func GetTimeDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		duration, err := time.ParseDuration(value)
		if err != nil {
			log.Fatalf("bad value converting %s to time.Duration: %v", key, err)
		}
		return duration
	}
	return fallback
}
