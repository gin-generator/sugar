package bootstrap

import (
	"context"
	"fmt"
	"github.com/gin-generator/sugar/package/logger"
	"github.com/gin-generator/sugar/package/mysql"
	"github.com/gin-generator/sugar/package/pgsql"
	"github.com/gin-gonic/gin"
)

type ServerType int

const (
	ServerHttp ServerType = iota
	ServerWebsocket
	ServerGrpc
)

// Server interface
type Server interface {
	Run(cfg *Config)
	GetEngine() *gin.Engine
	Use(middleware ...gin.HandlerFunc)
}

type Option interface {
	apply(*Bootstrap)
}

type optionFunc func(*Bootstrap)

func (o optionFunc) apply(b *Bootstrap) {
	o(b)
}

type Bootstrap struct {
	context.Context

	// App config
	cfg *Config

	// Server instance
	server Server

	// Initializers
	initializers []Initializer
}

// Initializer interface
type Initializer interface {
	Init(cfg *Config) error
	Name() string
}

// loggerInitializer
type loggerInitializer struct{}

func (l *loggerInitializer) Init(cfg *Config) error {
	logger.NewLogger(cfg.Logger)
	return nil
}

func (l *loggerInitializer) Name() string {
	return "Logger"
}

// mysqlInitializer MySQL
type mysqlInitializer struct{}

func (m *mysqlInitializer) Init(cfg *Config) error {
	if len(cfg.Database.Mysql) > 0 {
		mysql.NewMysql(cfg.Database.Mysql)
	}
	return nil
}

func (m *mysqlInitializer) Name() string {
	return "MySQL"
}

// pgsqlInitializer PostgresSQL
type pgsqlInitializer struct{}

func (p *pgsqlInitializer) Init(cfg *Config) error {
	if len(cfg.Database.Pgsql) > 0 {
		pgsql.NewPgsql(cfg.Database.Pgsql)
	}
	return nil
}

func (p *pgsqlInitializer) Name() string {
	return "PostgresSQL"
}

// NewBootstrap
/**
 * @description: create a new bootstrap instance
 * @param {ServerType} serverType, server type: http, grpc
 */
func NewBootstrap(serverType ServerType, opts ...Option) *Bootstrap {
	b := &Bootstrap{
		Context: context.Background(),
		cfg:     NewConfig("env.yaml", "./etc"),
		initializers: []Initializer{
			&loggerInitializer{},
			&mysqlInitializer{},
			&pgsqlInitializer{},
		},
	}

	// create server instance
	b.server = createServer(serverType, b.cfg)

	for _, opt := range opts {
		opt.apply(b)
	}

	b.start()
	return b
}

// createServer
/**
 * @description: create a server instance based on the server type
 * @param: {ServerType} serverType, server type: http, grpc
 */
func createServer(serverType ServerType, cfg *Config) Server {
	switch serverType {
	case ServerHttp:
		return newHttp(cfg.App.Env)
	case ServerWebsocket:
		panic("unsupported server websocket type")
	case ServerGrpc:
		panic("unsupported server grpc type")
	default:
		panic("unsupported server type")
	}
}

// WithAppConfig
/**
 * @description: set app config
 * @param {Config} cfg
 */
func WithAppConfig(cfg *Config) Option {
	return optionFunc(func(b *Bootstrap) {
		b.cfg = cfg
	})
}

// WithGinEngine
/**
 * @description: set gin engine
 */
func WithGinEngine(engine *gin.Engine) Option {
	return optionFunc(func(b *Bootstrap) {
		if httpServer, ok := b.server.(*Http); ok {
			httpServer.Engine = engine
		}
	})
}

// WithHttpMiddleware
/**
 * @description: set http middleware
 * @param {...gin.HandlerFunc} middleware
 */
func WithHttpMiddleware(middleware ...gin.HandlerFunc) Option {
	return optionFunc(func(b *Bootstrap) {
		b.server.Use(middleware...)
	})
}

// WithHttpRouter
/**
 * @description: set http router
 */
func WithHttpRouter(registerRouter RegisterRouter) Option {
	return optionFunc(func(b *Bootstrap) {
		registerRouter(b.server.GetEngine())
	})
}

// WithInitializers
/**
 * @description: set initializers
 * @param {...Initializer} initializers
 */
func WithInitializers(initializers ...Initializer) Option {
	return optionFunc(func(b *Bootstrap) {
		b.initializers = append(b.initializers, initializers...)
	})
}

// start
/**
 * @description: 执行所有初始化器
 */
func (b *Bootstrap) start() {
	for _, initializer := range b.initializers {
		if err := initializer.Init(b.cfg); err != nil {
			panic(fmt.Sprintf("Failed to initialize %s: %v", initializer.Name(), err))
		}
	}
}

// Run
/**
 * @description: start server
 */
func (b *Bootstrap) Run() {
	b.server.Run(b.cfg)
}
