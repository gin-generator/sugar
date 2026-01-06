package database

import (
	"fmt"
	l "github.com/gin-generator/logger"
	_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// MysqlConfig MySQL configuration with validation tags
type MysqlConfig struct {
	Host               string        `validate:"required,ip"`
	Port               int           `validate:"required,gt=0,lte=65535"`
	Database           string        `validate:"omitempty"`
	Username           string        `validate:"required"`
	Password           string        `validate:"required"`
	Charset            string        `validate:"required"`
	ParseTime          bool          `validate:"required"`
	MultiStatements    bool          `validate:"required"`
	Loc                string        `validate:"required"`
	MaxIdleConnections *int          `validate:"omitempty,gte=0"`
	MaxOpenConnections *int          `validate:"omitempty,gt=0"`
	MaxLifeSeconds     *int          `validate:"omitempty,gt=0"`
	SkipVersion        bool          `validate:"omitempty,boolean"`
	Logger             *LoggerConfig `validate:"omitempty"`
}

// LoggerConfig logger configuration for database with validation tags
type LoggerConfig struct {
	Level         string `validate:"required,oneof=debug error warn info"`
	MaxSize       int    `validate:"required,gt=0"`
	MaxBackup     int    `validate:"required,gte=0"`
	MaxAge        int    `validate:"required,gt=0"`
	Compress      *bool  `validate:"omitempty"`
	LocalTime     *bool  `validate:"omitempty"`
	SlowThreshold int    `validate:"required,gt=0"` // milliseconds
}

// NewMysqlConnection creates a MySQL connection
func NewMysqlConnection(name string, cfg MysqlConfig) (*gorm.DB, error) {
	// Use provided database name or fallback to name parameter
	dbName := cfg.Database
	if dbName == "" {
		dbName = name
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&multiStatements=%v&loc=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		dbName,
		cfg.Charset,
		cfg.ParseTime,
		cfg.MultiStatements,
		cfg.Loc,
	)

	dbConfig := _mysql.New(_mysql.Config{
		DSN:                       dsn,
		SkipInitializeWithVersion: cfg.SkipVersion,
	})

	var gormConfig *gorm.Config
	if cfg.Logger != nil {
		compress := false
		if cfg.Logger.Compress != nil {
			compress = *cfg.Logger.Compress
		}

		localTime := false
		if cfg.Logger.LocalTime != nil {
			localTime = *cfg.Logger.LocalTime
		}

		logger := l.NewLogger(
			l.WithFileName(fmt.Sprintf("storage/logs/%s.log", name)),
			l.WithLevel(cfg.Logger.Level),
			l.WithTimeZone(localTime),
			l.WithMaxSize(cfg.Logger.MaxSize),
			l.WithMaxBackup(cfg.Logger.MaxBackup),
			l.WithMaxAge(cfg.Logger.MaxAge),
			l.WithCompress(compress),
		)
		_logger := l.NewGormLogger(logger, l.WithSlowThreshold(time.Duration(cfg.Logger.SlowThreshold)*time.Millisecond))
		gormConfig = &gorm.Config{Logger: _logger}
	} else {
		gormConfig = &gorm.Config{}
	}

	db, err := gorm.Open(dbConfig, gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if cfg.MaxOpenConnections != nil && *cfg.MaxOpenConnections > 0 {
		sqlDB.SetMaxOpenConns(*cfg.MaxOpenConnections)
	}

	if cfg.MaxIdleConnections != nil && *cfg.MaxIdleConnections > 0 {
		sqlDB.SetMaxIdleConns(*cfg.MaxIdleConnections)
	}

	if cfg.MaxLifeSeconds != nil && *cfg.MaxLifeSeconds > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(*cfg.MaxLifeSeconds) * time.Second)
	}

	return db, nil
}
