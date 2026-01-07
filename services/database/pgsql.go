package database

import (
	"fmt"
	l "github.com/gin-generator/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

// PgsqlConfig PostgresSQL configuration with validation tags
type PgsqlConfig struct {
	Host                 string        `validate:"required,ip"`
	Port                 int           `validate:"required,gt=0,lte=65535"`
	Database             string        `validate:"required"`
	Username             string        `validate:"required"`
	Password             string        `validate:"required"`
	Timezone             string        `validate:"required"`
	PreferSimpleProtocol bool          `validate:"required"`
	MaxIdleConnections   *int          `validate:"omitempty,gte=0"`
	MaxOpenConnections   *int          `validate:"omitempty,gt=0"`
	MaxLifeSeconds       *int          `validate:"omitempty,gt=0"`
	Logger               *LoggerConfig `validate:"omitempty"`
}

// NewPgsqlConnection creates a PostgresSQL connection
func NewPgsqlConnection(name string, cfg PgsqlConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		cfg.Host,
		cfg.Username,
		cfg.Password,
		cfg.Database,
		cfg.Port,
		cfg.Timezone,
	)

	dbConfig := postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: cfg.PreferSimpleProtocol,
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
