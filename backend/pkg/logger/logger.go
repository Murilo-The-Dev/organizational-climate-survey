// Package logger provides a structured logging interface with configurable levels,
// context support, and field injection for observability in production systems.
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

// Level represents the severity level of a log message
type Level int

// Log levels in ascending order of severity
const (
	DebugLevel Level = iota // Detailed information for debugging
	InfoLevel              // General operational messages
	WarnLevel              // Warning messages for unexpected conditions
	ErrorLevel             // Error messages for failures
	FatalLevel             // Critical errors that cause program termination
)

// levelNames provides string representation for each log level
var levelNames = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}

// String returns the string representation of a log level
func (l Level) String() string {
	if l >= DebugLevel && l <= FatalLevel {
		return levelNames[l]
	}
	return "UNKNOWN"
}

// ParseLevel converts a string to its corresponding Level
// Returns InfoLevel for unrecognized strings as safe default
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

// Logger defines the interface for structured logging operations
// Supports field injection, context propagation, and multiple log levels
type Logger interface {
	Debug(msg string, args ...interface{}) // Log debug-level messages
	Info(msg string, args ...interface{})  // Log informational messages
	Warn(msg string, args ...interface{})  // Log warning messages
	Error(msg string, args ...interface{}) // Log error messages
	Fatal(msg string, args ...interface{}) // Log fatal messages and exit
	WithFields(fields map[string]interface{}) Logger // Add structured fields
	WithContext(ctx context.Context) Logger           // Add context for tracing
}

// Config holds logger configuration parameters
type Config struct {
	Level      Level     // Minimum level to log
	Output     io.Writer // Destination for log output
	TimeFormat string    // Time format for timestamps
	AddCaller  bool      // Whether to include caller information
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() *Config {
	return &Config{
		Level:      InfoLevel,
		Output:     os.Stdout,
		TimeFormat: time.RFC3339, // ISO 8601 format for machine readability
		AddCaller:  true,         // Include file:line for debugging
	}
}

// logger is the concrete implementation of the Logger interface
type logger struct {
	config *Config                 // Logger configuration
	fields map[string]interface{} // Structured fields to include in logs
	ctx    context.Context        // Context for request tracing
}

// New creates a new Logger instance with the provided configuration
// Uses DefaultConfig() if config is nil
func New(config *Config) Logger {
	if config == nil {
		config = DefaultConfig()
	}
	return &logger{
		config: config,
		fields: make(map[string]interface{}),
	}
}

// shouldLog determines if a message should be logged based on level filtering
func (l *logger) shouldLog(level Level) bool {
	return level >= l.config.Level
}

// log is the core logging method that formats and outputs log messages
// Handles level filtering, timestamp formatting, caller info, and structured fields
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

// getFileName extracts just the filename from a full path
// Used to keep caller info concise
func getFileName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return path
}

// Level-specific logging methods for convenience and type safety

// Debug logs debug-level messages, typically for detailed troubleshooting
func (l *logger) Debug(msg string, args ...interface{}) {
	l.log(DebugLevel, msg, args...)
}

// Info logs informational messages for normal operation tracking
func (l *logger) Info(msg string, args ...interface{}) {
	l.log(InfoLevel, msg, args...)
}

// Warn logs warning messages for unexpected but recoverable conditions
func (l *logger) Warn(msg string, args ...interface{}) {
	l.log(WarnLevel, msg, args...)
}

// Error logs error messages for failures that don't terminate the program
func (l *logger) Error(msg string, args ...interface{}) {
	l.log(ErrorLevel, msg, args...)
}

// Fatal logs critical error messages and terminates the program
func (l *logger) Fatal(msg string, args ...interface{}) {
	l.log(FatalLevel, msg, args...)
}

// WithFields creates a new logger instance with additional structured fields
// Fields are merged with existing fields, with new fields taking precedence
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

// WithContext creates a new logger instance with attached context
// Useful for request tracing and distributed logging correlation
func (l *logger) WithContext(ctx context.Context) Logger {
	return &logger{
		config: l.config,
		fields: l.fields,
		ctx:    ctx,
	}
}

// NoopLogger is a logger implementation that discards all log messages
// Useful for testing or when logging needs to be disabled
type NoopLogger struct{}

func (n NoopLogger) Debug(msg string, args ...interface{})                     {}
func (n NoopLogger) Info(msg string, args ...interface{})                      {}
func (n NoopLogger) Warn(msg string, args ...interface{})                      {}
func (n NoopLogger) Error(msg string, args ...interface{})                     {}
func (n NoopLogger) Fatal(msg string, args ...interface{})                     {}
func (n NoopLogger) WithFields(fields map[string]interface{}) Logger           { return n }
func (n NoopLogger) WithContext(ctx context.Context) Logger                    { return n }