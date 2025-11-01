package bootstrap

import (
	"fmt"
	"github.com/gin-generator/sugar/package/logger"
	"github.com/gin-generator/sugar/package/mysql"
	"github.com/gin-generator/sugar/package/pgsql"
	"github.com/gin-generator/sugar/package/validator"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type Mode string

const (
	ModeDebug   Mode = "debug"
	ModeRelease Mode = "release"
	ModeTest    Mode = "test"
)

type App struct {
	Name string `validate:"required"`
	Host string `validate:"required,ip"`
	Port int    `validate:"required,gt=0,lte=65535"`
	Env  Mode   `validate:"required,oneof=debug release test"`
}

type Database struct {
	Mysql map[string]mysql.Mysql `validate:"dive,keys,required,endkeys,required"`
	Pgsql map[string]pgsql.Pgsql `validate:"omitzero,dive,keys,required,endkeys,required"`
}

type Config struct {
	App      *App           `validate:"required"`
	Logger   *logger.Logger `validate:"required"`
	Database *Database      `validate:"required,dive"`
}

// NewConfig
/**
 * @description: load config file
 * @param {string} filename
 * @param {string} path
 */
func NewConfig(filename, path string) (cfg *Config) {
	v := viper.New()
	v.SetConfigName(filename)

	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	v.SetConfigType(ext)
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// v.WatchConfig()

	cfg = new(Config)
	err = v.Unmarshal(cfg)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// check config
	err = validator.ValidateStruct(*cfg)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return
}
