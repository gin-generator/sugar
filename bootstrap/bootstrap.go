package bootstrap

import (
	"github.com/gin-generator/sugar/config"
	"github.com/gin-generator/sugar/foundation"
	"github.com/gin-generator/sugar/providers"
	"github.com/gin-gonic/gin"
)

type ServerType int

const (
	ServerHttp ServerType = iota
	ServerWebsocket
	ServerGrpc
)

// Server 服务器接口
type Server interface {
	Run(app *foundation.Application)
}

type Option interface {
	apply(*Bootstrap)
}

type optionFunc func(*Bootstrap)

func (o optionFunc) apply(b *Bootstrap) {
	o(b)
}

// Bootstrap application bootstrap
type Bootstrap struct {
	// Application container
	app *foundation.Application

	// Config manager
	cfg *config.Config

	// Server instance
	server Server
}

// NewBootstrap creates a new bootstrap instance
func NewBootstrap(serverType ServerType, opts ...Option) *Bootstrap {
	// Create application container
	app := foundation.NewApplication()

	// Load configuration
	cfg := config.NewConfig("env.yaml", "./etc")

	b := &Bootstrap{
		app: app,
		cfg: cfg,
	}

	// Store configuration in application container
	var appCfg map[string]interface{}
	err := cfg.UnmarshalKey("app", &appCfg)
	if err != nil {
		return nil
	}
	app.SetConfig("app", appCfg)

	// Register service providers
	b.registerProviders()

	// Boot service providers
	if err = app.Boot(); err != nil {
		panic(err)
	}

	// Create server instance
	b.server = createServer(serverType, app, cfg)

	// Apply options
	for _, opt := range opts {
		opt.apply(b)
	}

	return b
}

// registerProviders registers service providers
func (b *Bootstrap) registerProviders() {
	// Register core service providers
	b.app.Register(providers.NewLoggerServiceProvider(b.cfg))
	b.app.Register(providers.NewDatabaseServiceProvider(b.cfg))
	b.app.Register(providers.NewCacheServiceProvider(b.cfg))
	b.app.Register(providers.NewStorageServiceProvider(b.cfg))
	b.app.Register(providers.NewQueueServiceProvider(b.cfg))
}

// createServer creates a server instance based on the server type
func createServer(serverType ServerType, app *foundation.Application, cfg *config.Config) Server {
	env := cfg.GetString("app.env")

	switch serverType {
	case ServerHttp:
		return newHttp(env)
	case ServerWebsocket:
		// 可以在这里实现 WebSocket 服务器
		panic("websocket server not implemented yet")
	case ServerGrpc:
		return newGrpc()
	default:
		panic("unsupported server type")
	}
}

// WithConfig sets the config manager
func WithConfig(cfg *config.Config) Option {
	return optionFunc(func(b *Bootstrap) {
		b.cfg = cfg
	})
}

// WithProvider registers an additional service provider
func WithProvider(provider foundation.ServiceProvider) Option {
	return optionFunc(func(b *Bootstrap) {
		b.app.Register(provider)
	})
}

// WithGinEngine 设置 gin 引擎（仅适用于 HTTP 服务器）
func WithGinEngine(engine *gin.Engine) Option {
	return optionFunc(func(b *Bootstrap) {
		if httpServer, ok := b.server.(*Http); ok {
			httpServer.Engine = engine
		}
	})
}

// WithHttpMiddleware 设置 HTTP 中间件（仅适用于 HTTP 服务器）
func WithHttpMiddleware(middleware ...gin.HandlerFunc) Option {
	return optionFunc(func(b *Bootstrap) {
		if httpServer, ok := b.server.(*Http); ok {
			httpServer.Use(middleware...)
		}
	})
}

// WithHttpRouter 设置 HTTP 路由（仅适用于 HTTP 服务器）
func WithHttpRouter(registerRouter RegisterRouter) Option {
	return optionFunc(func(b *Bootstrap) {
		if httpServer, ok := b.server.(*Http); ok {
			registerRouter(httpServer.Engine)
		}
	})
}

// WithGrpcService 注册 gRPC 服务（仅适用于 gRPC 服务器）
func WithGrpcService(registerService RegisterGrpcService) Option {
	return optionFunc(func(b *Bootstrap) {
		if grpcServer, ok := b.server.(*Grpc); ok {
			registerService(grpcServer.Server)
		}
	})
}

// Run starts the server
func (b *Bootstrap) Run() {
	b.server.Run(b.app)
}

// App returns the application container
func (b *Bootstrap) App() *foundation.Application {
	return b.app
}
