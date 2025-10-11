package bootstrap

import (
	"context"
	"fmt"
	"github.com/gin-generator/sugar/middleware"
	"github.com/gin-generator/sugar/package/logger"
	"github.com/gin-generator/sugar/package/mysql"
	"github.com/gin-gonic/gin"
)

type Engine interface {
	*gin.Engine | *GrpcEngine
}

type Server string

const (
	ServerHttp Server = "http"
	ServerGrpc Server = "grpc"
)

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
	Config

	// Http server engine
	HttpEngine *gin.Engine

	// Grpc server engine
	GrpcEngine *GrpcEngine
}

// NewBootstrap
/**
 * @description: create a new http bootstrap instance
 * @param {Server} server, server type: http, grpc
 */
func NewBootstrap(server Server, opts ...Option) *Bootstrap {
	b := &Bootstrap{
		Context: context.Background(),
		Config:  NewConfig("env.yaml", "./etc"),
	}

	switch server {
	case ServerHttp:
		gin.SetMode(string(b.Config.App.Env))
		b.HttpEngine = gin.New()
		b.use()
	case ServerGrpc:
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
func WithAppConfig(cfg Config) Option {
	return optionFunc(func(b *Bootstrap) {
		b.Config = cfg
	})
}

// WithGinEngine
/**
 * @description: set gin engine
 */
func WithGinEngine(engine *gin.Engine) Option {
	return optionFunc(func(b *Bootstrap) {
		b.HttpEngine = engine
		b.use()
	})
}

// use
/**
 * @description: use global middleware
 */
func (b *Bootstrap) use() {
	b.HttpEngine.Use(
		middleware.Recovery(),
		middleware.Logger(),
		middleware.Cors(),
	)
}

// start
/**
 * @description: start base server, e.g. mysql,redis,cache etc.
 */
func (b *Bootstrap) start() {
	// Setup logger
	logger.NewLogger(b.Config.Logger)

	// Setup mysql
	if len(b.Config.Database.Mysql) > 0 {
		mysql.NewMysql(b.Config.Database.Mysql)
	}

	// Setup pgsql
}

// RunHttp
/**
 * @description: start server
 */
func (b *Bootstrap) RunHttp() {
	RegisterDemoApiRoute(b.HttpEngine)
	fmt.Println(fmt.Sprintf("%s serve start: %s:%d...",
		b.App.Name, b.App.Host, b.App.Port))
	err := b.HttpEngine.Run(fmt.Sprintf("%s:%d", b.App.Host, b.App.Port))
	if err != nil {
		panic("Unable to start server, error: " + err.Error())
	}
}
