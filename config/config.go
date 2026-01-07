package config

import (
	"fmt"
	"github.com/gin-generator/sugar/package/validator"
	"github.com/gin-generator/sugar/services/database"
	"github.com/gin-generator/sugar/services/logger"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

// Mode application mode
type Mode string

const (
	ModeDebug   Mode = "debug"
	ModeRelease Mode = "release"
	ModeTest    Mode = "test"
)

// App application configuration
type App struct {
	Name   string `validate:"required"`
	Server string `validate:"required,oneof=http grpc websocket"`
	Host   string `validate:"required"`
	Port   int    `validate:"required,gt=0,lte=65535"`
	Env    Mode   `validate:"required,oneof=debug release test"`
}

// Database database configuration for validation
type Database struct {
	Mysql map[string]database.MysqlConfig `validate:"omitempty,dive"`
	Pgsql map[string]database.PgsqlConfig `validate:"omitempty,dive"`
}

// Config configuration structure
type Config struct {
	App      App           `validate:"required"`
	Logger   logger.Config `validate:"required"`
	Database Database      `validate:"omitempty"`
}

// NewConfig creates and validates configuration from file
func NewConfig(filename, path string) *Config {
	v := viper.New()
	v.SetConfigName(filename)

	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	v.SetConfigType(ext)
	v.AddConfigPath(path)

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error reading config file: %s", err))
	}

	// Parse configuration
	config := new(Config)
	err = v.Unmarshal(config)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshaling config: %s", err))
	}

	// Validate configuration
	err = validator.ValidateStruct(*config)
	if err != nil {
		panic(fmt.Errorf("fatal error validating config: %s", err))
	}

	return config
}
