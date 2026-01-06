package pgsql

import (
	"fmt"
	l "github.com/gin-generator/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Pgsql struct {
	Host                 string  `validate:"required,ip"`
	Port                 int     `validate:"required,gt=0,lte=65535"`
	Database             string  `validate:"required"`
	Username             string  `validate:"required"`
	Password             string  `validate:"required"`
	Timezone             string  `validate:"required"`
	PreferSimpleProtocol bool    `validate:"required"`
	MaxIdleConnections   *int    `validate:"omitempty,gte=0"`
	MaxOpenConnections   *int    `validate:"omitempty,gt=0"`
	MaxLifeSeconds       *int    `validate:"omitempty,gt=0"`
	Logger               *Logger `validate:"omitempty"`
}

type Logger struct {
	Level         string `validate:"required,oneof=debug error warn info"`
	MaxSize       int    `validate:"required,gt=0"`
	MaxBackup     int    `validate:"required,gte=0"`
	MaxAge        int    `validate:"required,gt=0"`
	Compress      *bool  `validate:"omitempty"`
	LocalTime     *bool  `validate:"omitempty"`
	SlowThreshold int    `validate:"required,gt=0"` // in milliseconds
}

type Database struct {
	Default *gorm.DB
	db      map[string]*gorm.DB
	mu      sync.RWMutex
}

var PG Database

// NewPgsql
/**
 * @description: initialize pgsql connections
 */
func NewPgsql(databases map[string]Pgsql) {
	PG.db = make(map[string]*gorm.DB)
	var firstDB *gorm.DB

	for dbname, v := range databases {
		db, err := connect(dbname, v)
		if err != nil {
			continue
		}

		PG.mu.Lock()
		PG.db[dbname] = db
		if firstDB == nil {
			firstDB = db
		}
		PG.mu.Unlock()
	}

	if firstDB != nil {
		PG.Default = firstDB
	}
}

// connect
/**
 * @description: connect pgsql
 */
func connect(dbname string, cfg Pgsql) (db *gorm.DB, err error) {
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
			l.WithFileName(fmt.Sprintf("storage/logs/%s.log", dbname)),
			l.WithLevel(cfg.Logger.Level),
			l.WithTimeZone(localTime),
			l.WithMaxSize(cfg.Logger.MaxSize),
			l.WithMaxBackup(cfg.Logger.MaxBackup),
			l.WithMaxAge(cfg.Logger.MaxAge),
			l.WithCompress(compress),
		)
		_logger := l.NewGormLogger(logger, l.WithSlowThreshold(time.Duration(cfg.Logger.SlowThreshold)*time.Millisecond))
		gormConfig = &gorm.Config{
			Logger: _logger,
		}
	} else {
		gormConfig = &gorm.Config{}
	}

	db, err = gorm.Open(dbConfig, gormConfig)
	if err != nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	if cfg.MaxOpenConnections != nil {
		sqlDB.SetMaxOpenConns(*cfg.MaxOpenConnections)
	}

	if cfg.MaxIdleConnections != nil {
		sqlDB.SetMaxIdleConns(*cfg.MaxIdleConnections)
	}

	if cfg.MaxLifeSeconds != nil {
		sqlDB.SetConnMaxLifetime(time.Duration(*cfg.MaxLifeSeconds) * time.Second)
	}
	return
}

// Use
/**
 * @description: get pgsql connection by dbname
 */
func (d *Database) Use(dbname string) *gorm.DB {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if db, ok := d.db[dbname]; ok {
		return db
	}
	return nil
}
