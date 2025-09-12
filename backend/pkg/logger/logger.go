package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

var ( 
	logLevelNames = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

func (l LogLevel) String() string {
	if l >= DEBUG && l <= FATAL {
		return logLevelNames[l]
	}
	return "UNKNOWN"
}

type Logger struct {
	mu       sync.Mutex
	minLevel LogLevel
	logger   *log.Logger
}

func NewLogger(levelStr string) *Logger {
	minLevel := INFO
	switch levelStr {
	case "debug":
		minLevel = DEBUG
	case "info":
		minLevel = INFO
	case "warn":
		minLevel = WARN
	case "error":
		minLevel = ERROR
	case "fatal":
		minLevel = FATAL
	}

	return &Logger{
		minLevel: minLevel,
		logger:   log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *Logger) log(level LogLevel, format string, v ...interface{}) {
	if level < l.minLevel {
		return
	}
	
	l.mu.Lock()
	defer l.mu.Unlock()
	
	prefix := fmt.Sprintf("[%s] ", level.String())
	l.logger.Output(3, prefix+fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(DEBUG, format, v...)
}

func (l *Logger) Info(format string, v ...interface{}) {
	l.log(INFO, format, v...)
}

func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(WARN, format, v...)
}

func (l *Logger) Error(format string, v ...interface{}) {
	l.log(ERROR, format, v...)
}

func (l *Logger) Fatal(format string, v ...interface{}) {
	l.log(FATAL, format, v...)
	os.Exit(1)
}

// Exemplo de uso global (opcional, mas comum)
var defaultLogger = NewLogger("info")

func Debug(format string, v ...interface{}) { defaultLogger.Debug(format, v...) }
func Info(format string, v ...interface{})  { defaultLogger.Info(format, v...) }
func Warn(format string, v ...interface{})  { defaultLogger.Warn(format, v...) }
func Error(format string, v ...interface{}) { defaultLogger.Error(format, v...) }
func Fatal(format string, v ...interface{}) { defaultLogger.Fatal(format, v...) }

// SetLevel permite alterar o nível de log em runtime
func SetLevel(levelStr string) {
	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()
	switch levelStr {
	case "debug":
		defaultLogger.minLevel = DEBUG
	case "info":
		defaultLogger.minLevel = INFO
	case "warn":
		defaultLogger.minLevel = WARN
	case "error":
		defaultLogger.minLevel = ERROR
	case "fatal":
		defaultLogger.minLevel = FATAL
	default:
		defaultLogger.minLevel = INFO
	}
}