package bootstrap

import (
	"context"
	"github.com/gin-generator/sugar/middleware"
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
		b.HttpEngine = gin.New()
		b.use()
	case ServerGrpc:
	default:
		panic("unsupported server type")
	}

	for _, opt := range opts {
		opt.apply(b)
	}

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
 * @param {Config} cfg
 */
func (b *Bootstrap) start(cfg Config) {

}

// RunHttp
/**
 * @description: start server
 */
func (b *Bootstrap) RunHttp() {

}
