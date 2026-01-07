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

// Server server interface
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

	// Server instance
	server Server
}

// NewBootstrap creates a new bootstrap instance
func NewBootstrap(serverType ServerType, opts ...Option) *Bootstrap {
	// Create application container
	app := foundation.NewApplication()

	// Load and store configuration in application
	app.Config = config.NewConfig("env.yaml", "./etc")

	b := &Bootstrap{
		app:    app,
		server: nil,
	}

	// Register service providers
	b.registerProviders()

	// Boot service providers
	if err := app.Boot(); err != nil {
		panic(err)
	}

	// Create server instance
        b.server = createServer(serverType, app)

	// Apply options
	for _, opt := range opts {
		opt.apply(b)
	}

	return b
}

// registerProviders registers service providers
func (b *Bootstrap) registerProviders() {
	// Register core service providers
	b.app.Register(providers.NewLoggerServiceProvider())
	b.app.Register(providers.NewDatabaseServiceProvider())
	b.app.Register(providers.NewCacheServiceProvider())
	b.app.Register(providers.NewStorageServiceProvider())
	b.app.Register(providers.NewQueueServiceProvider())
}

// createServer creates a server instance based on the server type
func createServer(serverType ServerType, app *foundation.Application) Server {
	env := string(app.Config.App.Env)

	switch serverType {
	case ServerHttp:
		return newHttp(env)
	case ServerWebsocket:
		panic("websocket server not implemented yet")
	case ServerGrpc:
		return newGrpc()
	default:
		panic("unsupported server type")
	}
}

// WithConfig sets the config (deprecated, config is auto-loaded)
func WithConfig(cfg *config.Config) Option {
	return optionFunc(func(b *Bootstrap) {
		b.app.Config = cfg
	})
}

// WithProvider registers an additional service provider
func WithProvider(provider foundation.ServiceProvider) Option {
	return optionFunc(func(b *Bootstrap) {
		b.app.Register(provider)
	})
}

// WithGinEngine sets the Gin engine (only for HTTP server)
func WithGinEngine(engine *gin.Engine) Option {
	return optionFunc(func(b *Bootstrap) {
		if httpServer, ok := b.server.(*Http); ok {
			httpServer.Engine = engine
		}
	})
}

// WithHttpMiddleware sets HTTP middleware (only for HTTP server)
func WithHttpMiddleware(middleware ...gin.HandlerFunc) Option {
	return optionFunc(func(b *Bootstrap) {
		if httpServer, ok := b.server.(*Http); ok {
			httpServer.Use(middleware...)
		}
	})
}

// WithHttpRouter sets HTTP routes (only for HTTP server)
func WithHttpRouter(registerRouter RegisterRouter) Option {
	return optionFunc(func(b *Bootstrap) {
		if httpServer, ok := b.server.(*Http); ok {
			registerRouter(httpServer.Engine)
		}
	})
}

// WithGrpcService sets gRPC services (only for gRPC server)
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
