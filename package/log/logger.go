package log

import "github.com/gin-generator/logger"

type Logger struct {
	Level     string
	Filename  string
	MaxSize   int
	MaxBackup int
	MaxAge    int
	Compress  bool
	LocalTime bool
}

var log *logger.Logger

func NewLogger(cfg Logger) {
	log = logger.NewLogger(
		logger.WithLevel(cfg.Level),
		logger.WithFileName(cfg.Filename),
		logger.WithMaxSize(cfg.MaxSize),
		logger.WithMaxBackup(cfg.MaxBackup),
		logger.WithMaxAge(cfg.MaxAge),
		logger.WithCompress(cfg.Compress),
		logger.WithTimeZone(cfg.LocalTime),
	)
}
