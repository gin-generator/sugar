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
	Name string `validate:"required"`
	Host string `validate:"required"`
	Port int    `validate:"required,gt=0,lte=65535"`
	Env  Mode   `validate:"required,oneof=debug release test"`
}

// Database database configuration for validation
type Database struct {
	Mysql map[string]database.MysqlConfig `validate:"omitempty,dive"`
	Pgsql map[string]database.PgsqlConfig `validate:"omitempty,dive"`
}

// GlobalConfig configuration structure for validation
type GlobalConfig struct {
	App      App           `validate:"required"`
	Logger   logger.Config `validate:"required"`
	Database Database      `validate:"required"`
}

// Config configuration manager
type Config struct {
	v *viper.Viper
}

// NewConfig creates a configuration manager and validates required fields
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

	// Validate configuration structure
	config := new(GlobalConfig)
	err = v.Unmarshal(config)
	if err != nil {
		panic(fmt.Errorf("fatal error unmarshaling config: %s", err))
	}

	// Validate required fields
	err = validator.ValidateStruct(*config)
	if err != nil {
		panic(fmt.Errorf("fatal error validating config: %s", err))
	}

	return &Config{v: v}
}

// Get retrieves a configuration value
func (c *Config) Get(key string) interface{} {
	return c.v.Get(key)
}

// GetString retrieves a string configuration value
func (c *Config) GetString(key string) string {
	return c.v.GetString(key)
}

// GetInt retrieves an integer configuration value
func (c *Config) GetInt(key string) int {
	return c.v.GetInt(key)
}

// GetBool retrieves a boolean configuration value
func (c *Config) GetBool(key string) bool {
	return c.v.GetBool(key)
}

// GetStringMap retrieves a map configuration value
func (c *Config) GetStringMap(key string) map[string]interface{} {
	return c.v.GetStringMap(key)
}

// Unmarshal parses configuration into a struct
func (c *Config) Unmarshal(rawVal interface{}) error {
	return c.v.Unmarshal(rawVal)
}

// UnmarshalKey parses a specific key's configuration into a struct
func (c *Config) UnmarshalKey(key string, rawVal interface{}) error {
	return c.v.UnmarshalKey(key, rawVal)
}
