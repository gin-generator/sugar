package bootstrap

import (
	"context"
	"github.com/gin-generator/sugar/package/logger"
	"github.com/gin-generator/sugar/package/mysql"
	"github.com/gin-gonic/gin"
)

type Server string

const (
	ServerHttp      Server = "http"
	ServerWebsocket Server = "websocket"
	ServerGrpc      Server = "grpc"
)

type Option interface {
	apply(*Bootstrap)
}

type optionFunc func(*Bootstrap)

func (o optionFunc) apply(b *Bootstrap) {
	o(b)
}

type Ability interface {
	Run()
}

type Bootstrap struct {
	context.Context

	// App config
	cfg *Config

	// Http server engine
	http *Http

	// Grpc server engine
	grpc *Grpc
}

// NewBootstrap
/**
 * @description: create a new http bootstrap instance
 * @param {Server} server, server type: http, grpc
 */
func NewBootstrap(server Server, opts ...Option) *Bootstrap {
	b := &Bootstrap{
		Context: context.Background(),
		cfg:     NewConfig("env.yaml", "./etc"),
	}

	switch server {
	case ServerHttp:
		b.http = newHttp(b.cfg.App.Env)
	case ServerWebsocket:
		panic("unsupported server websocket type")
	case ServerGrpc:
		panic("unsupported server grpc type")
	default:
		panic("unsupported server type")
	}

	for _, opt := range opts {
		opt.apply(b)
	}

	b.start()
	return b
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
		b.http.Engine = engine
	})
}

// WithHttpMiddleware
/**
 * @description: set http middleware
 * @param {...gin.HandlerFunc} middleware
 */
func WithHttpMiddleware(middleware ...gin.HandlerFunc) Option {
	return optionFunc(func(b *Bootstrap) {
		b.http.Use(middleware...)
	})
}

// WithHttpRouter
/**
 * @description: set http router
 */
func WithHttpRouter(registerRouter RegisterRouter) Option {
	return optionFunc(func(b *Bootstrap) {
		registerRouter(b.http.Engine)
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

	// Setup pgsql
}

// Run
/**
 * @description: start server
 */
func (b *Bootstrap) Run() {
	if b.http != nil {
		b.http.run(b.cfg)
	}
}
