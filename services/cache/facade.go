package cache

import (
	"context"
	"fmt"
	"time"
)

// Global cache manager instance
var manager *Manager

// SetManager sets the global cache manager
func SetManager(m *Manager) {
	manager = m
}

// Get retrieves a value from cache (Facade pattern)
func Get(ctx context.Context, key string) (string, error) {
	if manager == nil {
		return "", fmt.Errorf("cache manager not initialized")
	}
	cache, err := manager.Cache()
	if err != nil {
		return "", err
	}
	return cache.Get(ctx, key)
}

// Set stores a value in cache (Facade pattern)
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if manager == nil {
		return fmt.Errorf("cache manager not initialized")
	}
	cache, err := manager.Cache()
	if err != nil {
		return err
	}
	return cache.Set(ctx, key, value, expiration)
}

// Delete removes a value from cache (Facade pattern)
func Delete(ctx context.Context, key string) error {
	if manager == nil {
		return fmt.Errorf("cache manager not initialized")
	}
	cache, err := manager.Cache()
	if err != nil {
		return err
	}
	return cache.Delete(ctx, key)
}

// Store gets a cache store by name (Facade pattern)
func Store(name string) (Cache, error) {
	if manager == nil {
		return nil, fmt.Errorf("cache manager not initialized")
	}
	return manager.Store(name)
}
