// Package config fornece a estrutura e funções para carregar a configuração da aplicação a partir de variáveis de ambiente.
package config

import (
	"fmt"
	"os"
)

// Config agrupa todas as configurações da aplicação, incluindo App, Database, JWT e Log.
type Config struct {
	App struct {
		Name string // Nome da aplicação
		Port string // Porta em que a aplicação será executada
		Env  string // Ambiente (development, production, etc.)
	}
	Database struct {
		Host     string // Host do banco de dados
		Port     string // Porta do banco de dados
		User     string // Usuário do banco
		Password string // Senha do banco
		DBName   string // Nome do banco
		SSLMode  string // Modo SSL
	}
	JWT struct {
		Secret string // Chave secreta para JWT
	}
	Log struct {
		Level string // Nível de log (debug, info, etc.)
	}
}

// LoadConfig lê as variáveis de ambiente e preenche a struct Config, aplicando defaults quando necessário.
func LoadConfig() (*Config, error) {
	cfg := &Config{}

	cfg.App.Name = getEnvWithDefault("APP_NAME", "organizational-climate-survey")
	cfg.App.Port = getEnvWithDefault("APP_PORT", "8080")
	cfg.App.Env = getEnvWithDefault("APP_ENV", "development")

	cfg.Database.Host = getEnvWithDefault("DB_HOST", "localhost")
	cfg.Database.Port = getEnvWithDefault("DB_PORT", "5432")
	cfg.Database.User = getEnvWithDefault("DB_USER", "postgres")
	cfg.Database.Password = os.Getenv("DB_PASS")
	cfg.Database.DBName = getEnvWithDefault("DB_NAME", "organizational_climate")
	cfg.Database.SSLMode = getEnvWithDefault("DB_SSLMODE", "disable")

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	cfg.Log.Level = getEnvWithDefault("LOG_LEVEL", "debug")

	// Validações obrigatórias
	if cfg.Database.Password == "" {
		return nil, fmt.Errorf("DB_PASS não configurado nas variáveis de ambiente")
	}
	if cfg.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET não configurado nas variáveis de ambiente")
	}

	return cfg, nil
}

// getEnvWithDefault retorna o valor da variável de ambiente ou um valor padrão caso não esteja definida.
func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
