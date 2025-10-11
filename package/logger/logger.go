package logger

import _logger "github.com/gin-generator/logger"

type Logger struct {
	Level     string
	Filename  string
	MaxSize   int
	MaxBackup int
	MaxAge    int
	Compress  bool
	LocalTime bool
}

var Log *_logger.Logger

func NewLogger(cfg Logger) {

	Log = _logger.NewLogger(
		_logger.WithLevel(cfg.Level),
		_logger.WithFileName(cfg.Filename),
		_logger.WithMaxSize(cfg.MaxSize),
		_logger.WithMaxBackup(cfg.MaxBackup),
		_logger.WithMaxAge(cfg.MaxAge),
		_logger.WithCompress(cfg.Compress),
		_logger.WithTimeZone(cfg.LocalTime),
	)
}
