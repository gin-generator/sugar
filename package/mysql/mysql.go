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
	Host                     string
	Port                     int
	Username                 string
	Password                 string
	Charset                  string
	ParseTime                bool
	MultiStatements          bool
	Loc                      string
	MaxIdleConnections       int
	MaxOpenConnections       int
	MaxLifeSeconds           int
	KipInitializeWithVersion bool
	Logger
}

type Logger struct {
	Level         string
	MaxSize       int
	MaxBackup     int
	MaxAge        int
	Compress      bool
	LocalTime     bool
	SlowThreshold int
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
		SkipInitializeWithVersion: cfg.KipInitializeWithVersion,
	})

	logger := l.NewLogger(
		l.WithFileName(fmt.Sprintf("storage/logs/%s.log", dbname)),
		l.WithLevel(cfg.Level),
		l.WithTimeZone(cfg.LocalTime),
		l.WithMaxSize(cfg.MaxSize),
		l.WithMaxBackup(cfg.MaxBackup),
		l.WithMaxAge(cfg.MaxAge),
		l.WithCompress(cfg.Compress),
	)
	_logger := l.NewGormLogger(logger, l.WithSlowThreshold(time.Duration(cfg.SlowThreshold)*time.Millisecond))
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
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxLifeSeconds) * time.Second)
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
