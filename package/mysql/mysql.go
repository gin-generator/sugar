package mysql

import (
	"fmt"
	l "github.com/gin-generator/logger"
	_mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Mysql struct {
	Host               string  `validate:"required,ip"`
	Port               int     `validate:"required,gt=0,lte=65535"`
	Username           string  `validate:"required"`
	Password           string  `validate:"required"`
	Charset            string  `validate:"required"`
	ParseTime          bool    `validate:"required"`
	MultiStatements    bool    `validate:"required"`
	Loc                string  `validate:"required"`
	MaxIdleConnections *int    `validate:"omitempty,gte=0"`
	MaxOpenConnections *int    `validate:"omitempty,gt=0"`
	MaxLifeSeconds     *int    `validate:"omitempty,gt=0"`
	SkipVersion        bool    `validate:"omitempty,boolean"`
	Logger             *Logger `validate:"omitempty"`
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

var DB Database

// NewMysql
/**
 * @description: initialize mysql connections
 */
func NewMysql(databases map[string]Mysql) {
	DB.db = make(map[string]*gorm.DB)
	var firstDB *gorm.DB

	for dbname, v := range databases {
		db, err := connect(dbname, v)
		if err != nil {
			continue
		}

		DB.mu.Lock()
		DB.db[dbname] = db
		if firstDB == nil {
			firstDB = db
		}
		DB.mu.Unlock()
	}

	if firstDB != nil {
		DB.Default = firstDB
	}
}

// connect
/**
 * @description: connect mysql
 */
func connect(dbname string, cfg Mysql) (db *gorm.DB, err error) {
	var dbConfig gorm.Dialector
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&multiStatements=%v&loc=%v",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		dbname,
		cfg.Charset,
		cfg.ParseTime,
		cfg.MultiStatements,
		cfg.Loc,
	)
	dbConfig = _mysql.New(_mysql.Config{
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
 * @description: get mysql connection by dbname
 */
func (d *Database) Use(dbname string) *gorm.DB {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if db, ok := d.db[dbname]; ok {
		return db
	}
	return nil
}
