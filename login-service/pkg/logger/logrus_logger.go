// internal/adapters/logger/logrus_logger.go
package logruslogger

import "github.com/sirupsen/logrus"

// LogrusLogger is a concrete implementation of the Logger interface using logrus
type LogrusLogger struct {
	logger *logrus.Logger
}

// NewLogrusLogger creates a new Logrus-based logger
func NewLogrusLogger() *LogrusLogger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel) // Default level can be changed
	return &LogrusLogger{logger: logger}
}

// Info logs an info-level message
func (l *LogrusLogger) Info(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Info(msg)
}

// Error logs an error-level message
func (l *LogrusLogger) Error(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Error(msg)
}

// Warn logs a warning-level message
func (l *LogrusLogger) Warn(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Warn(msg)
}

// Debug logs a debug-level message
func (l *LogrusLogger) Debug(msg string, fields map[string]interface{}) {
	l.logger.WithFields(fields).Debug(msg)
}
