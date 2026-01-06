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
	Host string `validate:"required"`
	Port int    `validate:"required,gt=0,lte=65535"`
	Env  Mode   `validate:"required,oneof=debug release test"`
}

type Database struct {
	Mysql map[string]mysql.Mysql `validate:"omitempty,dive"`
	Pgsql map[string]pgsql.Pgsql `validate:"omitempty,dive"`
}

type Config struct {
	App      App           `validate:"required"`
	Logger   logger.Logger `validate:"required"`
	Database Database      `validate:"required"`
}

// NewConfig
/**
 * @description: load config file
 * @param {string} filename
 * @param {string} path
 */
func NewConfig(filename, path string) (config *Config) {
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

	config = new(Config)
	err = v.Unmarshal(config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// check config
	err = validator.ValidateStruct(*config)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return
}
