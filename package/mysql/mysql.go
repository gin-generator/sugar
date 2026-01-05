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
	Logger             *Logger `validate:"required"`
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
	db      sync.Map
}

var DB Database

// NewMysql
/**
 * @description: initialize mysql connections
 */
func NewMysql(databases map[string]Mysql) {
	var i int
	for dbname, v := range databases {
		i++
		_, ok := DB.db.Load(dbname)
		if !ok {
			db, err := connect(dbname, v)
			if err != nil {
				continue
			} else {
				DB.db.Store(dbname, db)
				if i == 1 {
					DB.Default = db
				}
			}
		}
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

	logger := l.NewLogger(
		l.WithFileName(fmt.Sprintf("storage/logs/%s.log", dbname)),
		l.WithLevel(cfg.Logger.Level),
		l.WithTimeZone(*cfg.Logger.LocalTime),
		l.WithMaxSize(cfg.Logger.MaxSize),
		l.WithMaxBackup(cfg.Logger.MaxBackup),
		l.WithMaxAge(cfg.Logger.MaxAge),
		l.WithCompress(*cfg.Logger.Compress),
	)
	_logger := l.NewGormLogger(logger, l.WithSlowThreshold(time.Duration(cfg.Logger.SlowThreshold)*time.Millisecond))
	db, err = gorm.Open(dbConfig, &gorm.Config{
		Logger: _logger,
	})
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
	if db, ok := d.db.Load(dbname); ok {
		return db.(*gorm.DB)
	}
	return nil
}
