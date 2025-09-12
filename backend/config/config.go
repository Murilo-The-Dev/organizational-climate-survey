package config

import (
	"fmt"
	"os"
)

type Config struct {
	App struct {
		Name string
		Port string
		Env  string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}
	JWT struct {
		Secret string
	}
	Log struct {
		Level string
	}
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.App.Name = getEnvWithDefault("APP_NAME", "organizational-climate-survey")
	cfg.App.Port = getEnvWithDefault("APP_PORT", "8080")
	cfg.App.Env = getEnvWithDefault("APP_ENV", "development")

	cfg.Database.Host = getEnvWithDefault("DB_HOST", "localhost")
	cfg.Database.Port = getEnvWithDefault("DB_PORT", "5432")
	cfg.Database.User = getEnvWithDefault("DB_USER", "postgres")
	cfg.Database.Password = os.Getenv("DB_PASS")
	cfg.Database.DBName = getEnvWithDefault("DB_NAME", "Atmos")
	cfg.Database.SSLMode = getEnvWithDefault("DB_SSLMODE", "disable")

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")

	cfg.Log.Level = getEnvWithDefault("LOG_LEVEL", "debug")

	// Validações básicas
	if cfg.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASS não configurado nas variáveis de ambiente")
	}
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET não configurado nas variáveis de ambiente")
	}

	return cfg, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}