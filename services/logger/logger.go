package logger

import (
	"fmt"
	_logger "github.com/gin-generator/logger"
)

// Config logger configuration with validation tags
type Config struct {
	Level     string `validate:"required,oneof=debug info warn error"`
	Filename  string `validate:"required"`
	MaxSize   int    `validate:"required,gt=0"`
	MaxBackup int    `validate:"required,gte=0"`
	MaxAge    int    `validate:"required,gt=0"`
	Compress  *bool  `validate:"omitempty"`
	LocalTime *bool  `validate:"omitempty"`
}

// Log global logger instance
var Log *_logger.Logger

// SetLogger sets the global logger instance
func SetLogger(log *_logger.Logger) {
	Log = log
}

// NewLoggerFromConfig creates a new logger instance from configuration
func NewLoggerFromConfig(cfg Config) *_logger.Logger {
	compress := false
	if cfg.Compress != nil {
		compress = *cfg.Compress
	}

	localTime := false
	if cfg.LocalTime != nil {
		localTime = *cfg.LocalTime
	}

	return _logger.NewLogger(
		_logger.WithLevel(cfg.Level),
		_logger.WithFileName(cfg.Filename),
		_logger.WithMaxSize(cfg.MaxSize),
		_logger.WithMaxBackup(cfg.MaxBackup),
		_logger.WithMaxAge(cfg.MaxAge),
		_logger.WithCompress(compress),
		_logger.WithTimeZone(localTime),
	)
}

// Info logs an Info level message (Facade pattern)
func Info(msg string, fields ...interface{}) error {
	if Log == nil {
		return fmt.Errorf("logger not initialized")
	}
	Log.Info(msg)
	return nil
}

// Debug logs a Debug level message (Facade pattern)
func Debug(msg string, fields ...interface{}) error {
	if Log == nil {
		return fmt.Errorf("logger not initialized")
	}
	Log.Debug(msg)
	return nil
}

// Warn logs a Warn level message (Facade pattern)
func Warn(msg string, fields ...interface{}) error {
	if Log == nil {
		return fmt.Errorf("logger not initialized")
	}
	Log.Warn(msg)
	return nil
}

// Error logs an Error level message (Facade pattern)
func Error(msg string, fields ...interface{}) error {
	if Log == nil {
		return fmt.Errorf("logger not initialized")
	}
	Log.Error(msg)
	return nil
}
