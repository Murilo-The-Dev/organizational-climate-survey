// Package config fornece a estrutura e funções para carregar a configuração da aplicação a partir de variáveis de ambiente.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config agrupa todas as configurações da aplicação, incluindo App, Database, JWT e Log.
type Config struct {
	App struct {
		Name                  string // Nome da aplicação
		Port                  string // Porta em que a aplicação será executada
		Env                   string // Ambiente (development, production, etc.)
		FrontendURL           string // Lista de origens permitidas no CORS (separadas por vírgula)
		RequestBodyLimitBytes int64  // Tamanho máximo do corpo de requisição
	}
	Database struct {
		Host                string // Host do banco de dados
		Port                string // Porta do banco de dados
		User                string // Usuário do banco
		Password            string // Senha do banco
		DBName              string // Nome do banco
		SSLMode             string // Modo SSL
		MaxOpenConns        int    // Máximo de conexões abertas
		MaxIdleConns        int    // Máximo de conexões ociosas
		ConnMaxLifetimeMins int    // Tempo de vida máximo de conexão em minutos
	}
	JWT struct {
		Secret string // Chave secreta para JWT
	}
	Crypto struct {
		HashSalt string // Salt para hashes de IP/fingerprint (submissões anônimas)
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
	cfg.App.FrontendURL = getEnvWithDefault("FRONTEND_URL", "http://localhost:3000")

	requestLimitBytes, err := getEnvAsInt64("REQUEST_BODY_LIMIT_BYTES", defaultRequestBodyLimitBytes)
	if err != nil {
		return nil, fmt.Errorf("REQUEST_BODY_LIMIT_BYTES inválido: %v", err)
	}
	cfg.App.RequestBodyLimitBytes = requestLimitBytes

	cfg.Database.Host = getEnvWithDefault("DB_HOST", "localhost")
	cfg.Database.Port = getEnvWithDefault("DB_PORT", "5432")
	cfg.Database.User = getEnvWithDefault("DB_USER", "postgres")
	cfg.Database.Password = os.Getenv("DB_PASS")
	cfg.Database.DBName = getEnvWithDefault("DB_NAME", "organizational_climate")
	cfg.Database.SSLMode = getEnvWithDefault("DB_SSLMODE", "disable")

	maxOpenConns, err := getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	if err != nil {
		return nil, fmt.Errorf("DB_MAX_OPEN_CONNS inválido: %v", err)
	}
	maxIdleConns, err := getEnvAsInt("DB_MAX_IDLE_CONNS", 5)
	if err != nil {
		return nil, fmt.Errorf("DB_MAX_IDLE_CONNS inválido: %v", err)
	}
	connLifetimeMins, err := getEnvAsInt("DB_CONN_MAX_LIFETIME_MINS", 5)
	if err != nil {
		return nil, fmt.Errorf("DB_CONN_MAX_LIFETIME_MINS inválido: %v", err)
	}

	cfg.Database.MaxOpenConns = maxOpenConns
	cfg.Database.MaxIdleConns = maxIdleConns
	cfg.Database.ConnMaxLifetimeMins = connLifetimeMins

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")

	cfg.Crypto.HashSalt = getEnvWithDefault("HASH_SALT", "default-salt-change-in-production-12345")

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

const defaultRequestBodyLimitBytes int64 = 1 << 20 // 1 MiB

func getEnvAsInt(key string, defaultValue int) (int, error) {
	value := stringsTrim(os.Getenv(key))
	if value == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func getEnvAsInt64(key string, defaultValue int64) (int64, error) {
	value := stringsTrim(os.Getenv(key))
	if value == "" {
		return defaultValue, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func stringsTrim(value string) string {
	return strings.TrimSpace(value)
}
