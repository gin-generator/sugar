package bootstrap

import (
	"fmt"
	"github.com/gin-generator/sugar/package/log"
	"github.com/gin-generator/sugar/package/mysql"
	"github.com/gin-generator/sugar/package/pgsql"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

type App struct {
	Name string
	Host string
	Port int
	Env  string
}

type Database struct {
	Mysql map[string]mysql.Mysql
	Pgsql map[string]pgsql.Pgsql
}

type Config struct {
	App
	log.Logger
	Database
}

func NewConfig(filename, path string) (cfg Config) {
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

	err = v.Unmarshal(&cfg)
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return
}
