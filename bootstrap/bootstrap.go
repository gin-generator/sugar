package bootstrap

import (
	"context"
	"github.com/gin-generator/sugar/package/logger"
	"github.com/gin-generator/sugar/package/mysql"
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

	// Server 实例（接口类型）
	server Server
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

// start
/**
 * @description: start base server, e.g. mysql,redis,cache etc.
 */
func (b *Bootstrap) start() {
	// Setup logger
	logger.NewLogger(b.cfg.Logger)

	// Setup mysql
	if len(b.cfg.Database.Mysql) > 0 {
		mysql.NewMysql(b.cfg.Database.Mysql)
	}

	// TODO：Setup pgsql
}

// Run
/**
 * @description: start server
 */
func (b *Bootstrap) Run() {
	b.server.Run(b.cfg)
}
