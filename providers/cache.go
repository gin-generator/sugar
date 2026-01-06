package providers

import (
	"fmt"
	"github.com/gin-generator/sugar/config"
	"github.com/gin-generator/sugar/foundation"
	"github.com/gin-generator/sugar/services/cache"
)

// CacheServiceProvider cache service provider
type CacheServiceProvider struct {
	cfg *config.Config
}

// NewCacheServiceProvider creates a cache service provider
func NewCacheServiceProvider(cfg *config.Config) *CacheServiceProvider {
	return &CacheServiceProvider{cfg: cfg}
}

// Register registers the service
func (p *CacheServiceProvider) Register(app *foundation.Application) {
	manager := cache.NewManager()
	app.Bind("cache", manager)
}

// Boot boots the service
func (p *CacheServiceProvider) Boot(app *foundation.Application) error {
	service, ok := app.Make("cache")
	if !ok {
		return fmt.Errorf("cache service not found")
	}

	manager := service.(*cache.Manager)

	// Get cache driver configuration
	driver := p.cfg.GetString("cache.drive")

	switch driver {
	case "redis":
		var redisCfg cache.RedisConfig
		if err := p.cfg.UnmarshalKey("cache.redis", &redisCfg); err != nil {
			return fmt.Errorf("failed to unmarshal redis config: %w", err)
		}

		redisStore := cache.NewRedisStore(redisCfg)
		manager.AddStore("redis", redisStore)
	default:
		// Add other cache drivers like memcached
	}

	// Set global Facade
	cache.SetManager(manager)

	return nil
}

// Name returns the service provider name
func (p *CacheServiceProvider) Name() string {
	return "Cache"
}
