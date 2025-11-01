package logger

import _logger "github.com/gin-generator/logger"

type Logger struct {
	Level     string `validate:"required,oneof=debug info warn error"`
	Filename  string `validate:"required"`
	MaxSize   int    `validate:"required,gt=0"`
	MaxBackup int    `validate:"required,gte=0"`
	MaxAge    int    `validate:"required,gt=0"`
	Compress  *bool  `validate:"omitempty"`
	LocalTime *bool  `validate:"omitempty"`
}

var Log *_logger.Logger

// NewLogger
/**
 * @description: create a new logger instance
 */
func NewLogger(cfg Logger) {
	Log = _logger.NewLogger(
		_logger.WithLevel(cfg.Level),
		_logger.WithFileName(cfg.Filename),
		_logger.WithMaxSize(cfg.MaxSize),
		_logger.WithMaxBackup(cfg.MaxBackup),
		_logger.WithMaxAge(cfg.MaxAge),
		_logger.WithCompress(*cfg.Compress),
		_logger.WithTimeZone(*cfg.LocalTime),
	)
}
