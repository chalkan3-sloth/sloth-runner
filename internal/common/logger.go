package common

import (
	"fmt"
	"log/slog"
	"os"
)

// Logger provides logging functionality
type Logger struct {
	logger *slog.Logger
}

// GetLogger returns a logger instance
func GetLogger() *Logger {
	return &Logger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

// Info logs an informational message
func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

// Success logs a success message
func (l *Logger) Success(msg string) {
	fmt.Printf("✅ %s\n", msg)
}

// Error logs an error message
func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}
