// Package logger implementa um sistema de logging estruturado com suporte a níveis,
// contexto e campos customizáveis para observabilidade em sistemas em produção.
package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

// Level representa o nível de severidade de uma mensagem de log
type Level int

// Níveis de log em ordem crescente de severidade
const (
	DebugLevel Level = iota // Informações detalhadas para depuração
	InfoLevel               // Mensagens operacionais gerais
	WarnLevel               // Avisos sobre condições inesperadas
	ErrorLevel              // Erros que não causam falha total
	FatalLevel              // Erros críticos que encerram o programa
)

// levelNames fornece a representação em texto para cada nível de log
var levelNames = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// String retorna a representação em texto do nível de log
func (l Level) String() string {
	if l >= DebugLevel && l <= FatalLevel {
		return levelNames[l]
	}
	return "UNKNOWN"
}

// ParseLevel converte uma string para seu nível correspondente
// Retorna InfoLevel para strings não reconhecidas como padrão seguro
func ParseLevel(s string) Level {
	switch strings.ToLower(s) {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn", "warning":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "fatal":
		return FatalLevel
	default:
		return InfoLevel
	}
}

// Logger define a interface para operações de logging estruturado
// Suporta injeção de campos, propagação de contexto e múltiplos níveis
type Logger interface {
	Debug(msg string, args ...interface{})           // Registra mensagens de depuração
	Info(msg string, args ...interface{})            // Registra mensagens informativas
	Warn(msg string, args ...interface{})            // Registra mensagens de aviso
	Error(msg string, args ...interface{})           // Registra mensagens de erro
	Fatal(msg string, args ...interface{})           // Registra erros fatais e encerra
	WithFields(fields map[string]interface{}) Logger // Adiciona campos estruturados
	WithContext(ctx context.Context) Logger          // Adiciona contexto para rastreamento
}

// Config armazena os parâmetros de configuração do logger
type Config struct {
	Level      Level     // Nível mínimo para logar
	Output     io.Writer // Destino para saída dos logs
	TimeFormat string    // Formato de data/hora para timestamps
	AddCaller  bool      // Se deve incluir informações do caller (arquivo:linha)
}

// DefaultConfig retorna uma configuração padrão sensata para o logger
func DefaultConfig() *Config {
	return &Config{
		Level:      InfoLevel,
		Output:     os.Stdout,
		TimeFormat: time.RFC3339, // Formato ISO 8601 para legibilidade por máquina
		AddCaller:  true,         // Inclui arquivo:linha para depuração
	}
}

// logger é a implementação concreta da interface Logger
type logger struct {
	config *Config                // Configuração do logger
	fields map[string]interface{} // Campos estruturados para incluir nos logs
	ctx    context.Context        // Contexto para rastreamento de requisições
}

// New cria uma nova instância do Logger com a configuração fornecida
// Usa DefaultConfig() se config for nil
func New(config *Config) Logger {
	if config == nil {
		config = DefaultConfig()
	}
	return &logger{
		config: config,
		fields: make(map[string]interface{}),
	}
}

// shouldLog determina se uma mensagem deve ser logada com base na filtragem de nível
func (l *logger) shouldLog(level Level) bool {
	return level >= l.config.Level
}

// log é o método central de logging que formata e envia mensagens
// Gerencia filtragem de nível, formatação de timestamp, info do caller e campos estruturados
func (l *logger) log(level Level, msg string, args ...interface{}) {
	if !l.shouldLog(level) {
		return
	}

	timestamp := time.Now().Format(l.config.TimeFormat)
	levelStr := level.String()

	// Format message with printf-style arguments
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	// Build structured log line
	var logLine strings.Builder
	logLine.WriteString(fmt.Sprintf("[%s] %s", levelStr, timestamp))

	// Add caller information for debugging (skips 3 frames to get actual caller)
	if l.config.AddCaller {
		if _, file, line, ok := runtime.Caller(3); ok {
			logLine.WriteString(fmt.Sprintf(" %s:%d", getFileName(file), line))
		}
	}

	// Add structured fields if present
	if len(l.fields) > 0 {
		logLine.WriteString(" fields=")
		logLine.WriteString(fmt.Sprintf("%+v", l.fields))
	}

	// Extract common context values for request tracing
	if l.ctx != nil {
		if requestID := l.ctx.Value("request_id"); requestID != nil {
			logLine.WriteString(fmt.Sprintf(" request_id=%v", requestID))
		}
	}

	logLine.WriteString(fmt.Sprintf(" msg=\"%s\"\n", msg))

	// Write to configured output destination
	fmt.Fprint(l.config.Output, logLine.String())

	// Fatal level triggers process termination
	if level == FatalLevel {
		os.Exit(1)
	}
}

// getFileName extrai apenas o nome do arquivo de um caminho completo
// Usado para manter informações do caller concisas
func getFileName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return path
}

// Métodos específicos para cada nível de log, para conveniência e type safety

// Debug registra mensagens de nível debug, tipicamente para troubleshooting detalhado
func (l *logger) Debug(msg string, args ...interface{}) {
	l.log(DebugLevel, msg, args...)
}

// Info registra mensagens informativas para acompanhamento da operação normal
func (l *logger) Info(msg string, args ...interface{}) {
	l.log(InfoLevel, msg, args...)
}

// Warn registra mensagens de aviso para condições inesperadas mas recuperáveis
func (l *logger) Warn(msg string, args ...interface{}) {
	l.log(WarnLevel, msg, args...)
}

// Error registra mensagens de erro para falhas que não terminam o programa
func (l *logger) Error(msg string, args ...interface{}) {
	l.log(ErrorLevel, msg, args...)
}

// Fatal registra mensagens de erro crítico e termina o programa
func (l *logger) Fatal(msg string, args ...interface{}) {
	l.log(FatalLevel, msg, args...)
}

// WithFields cria uma nova instância do logger com campos estruturados adicionais
// Os campos são mesclados com os existentes, com novos campos tendo precedência
func (l *logger) WithFields(fields map[string]interface{}) Logger {
	newFields := make(map[string]interface{})
	// Copy existing fields
	for k, v := range l.fields {
		newFields[k] = v
	}
	// Add/override with new fields
	for k, v := range fields {
		newFields[k] = v
	}

	return &logger{
		config: l.config,
		fields: newFields,
		ctx:    l.ctx,
	}
}

// WithContext cria uma nova instância do logger com contexto anexado
// Útil para rastreamento de requisições e correlação de logs distribuídos
func (l *logger) WithContext(ctx context.Context) Logger {
	return &logger{
		config: l.config,
		fields: l.fields,
		ctx:    ctx,
	}
}

// NoopLogger é uma implementação do logger que descarta todas as mensagens
// Útil para testes ou quando o logging precisa ser desabilitado
type NoopLogger struct{}

func (n NoopLogger) Debug(msg string, args ...interface{})           {}
func (n NoopLogger) Info(msg string, args ...interface{})            {}
func (n NoopLogger) Warn(msg string, args ...interface{})            {}
func (n NoopLogger) Error(msg string, args ...interface{})           {}
func (n NoopLogger) Fatal(msg string, args ...interface{})           {}
func (n NoopLogger) WithFields(fields map[string]interface{}) Logger { return n }
func (n NoopLogger) WithContext(ctx context.Context) Logger          { return n }
