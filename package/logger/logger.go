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
	// 设置默认值
	compress := false
	if cfg.Compress != nil {
		compress = *cfg.Compress
	}
	
	localTime := false
	if cfg.LocalTime != nil {
		localTime = *cfg.LocalTime
	}
	
	Log = _logger.NewLogger(
		_logger.WithLevel(cfg.Level),
		_logger.WithFileName(cfg.Filename),
		_logger.WithMaxSize(cfg.MaxSize),
		_logger.WithMaxBackup(cfg.MaxBackup),
		_logger.WithMaxAge(cfg.MaxAge),
		_logger.WithCompress(compress),
		_logger.WithTimeZone(localTime),
	)
}
