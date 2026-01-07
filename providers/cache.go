package providers

import (
	"github.com/gin-generator/sugar/foundation"
	"github.com/gin-generator/sugar/services/cache"
)

// CacheServiceProvider cache service provider
type CacheServiceProvider struct{}

// NewCacheServiceProvider creates a cache service provider
func NewCacheServiceProvider() *CacheServiceProvider {
	return &CacheServiceProvider{}
}

// Register registers the service
func (p *CacheServiceProvider) Register(app *foundation.Application) {
	manager := cache.NewManager()
	app.Bind(ServiceCache, manager)
}

// Boot boots the service
func (p *CacheServiceProvider) Boot(app *foundation.Application) error {
	manager := foundation.MustMake[*cache.Manager](app, ServiceCache)

	// TODO: 从 app.Config 中获取 cache 配置
	// 目前 Config 中没有定义 cache 配置，需要后续添加
	// cfg := app.Config.(*config.Config)
	// driver := cfg.Cache.Drive
	// switch driver {
	// case "redis":
	// redisStore := cache.NewRedisStore(cfg.Cache.Redis)
	// manager.AddStore("redis", redisStore)
	// default:
	// // Add other cache drivers like memcached
	// }

	// Set global Facade
	cache.SetManager(manager)

	return nil
}

// Name returns the service provider name
func (p *CacheServiceProvider) Name() string {
	return "Cache"
}
