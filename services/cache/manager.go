package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Cache interface
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Delete(ctx context.Context, key string) error
	Has(ctx context.Context, key string) (bool, error)
	Flush(ctx context.Context) error
}

// Manager cache manager
type Manager struct {
	stores       map[string]Cache
	mu           sync.RWMutex
	defaultStore string
}

// NewManager creates a new cache manager
func NewManager() *Manager {
	return &Manager{
		stores: make(map[string]Cache),
	}
}

// AddStore adds a cache store
func (m *Manager) AddStore(name string, store Cache) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stores[name] = store

	if m.defaultStore == "" {
		m.defaultStore = name
	}
}

// Store gets a cache store by name
func (m *Manager) Store(name string) (Cache, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	store, ok := m.stores[name]
	if !ok {
		return nil, fmt.Errorf("cache store %s not found", name)
	}
	return store, nil
}

// Cache gets the default cache store
func (m *Manager) Cache() (Cache, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.defaultStore == "" {
		return nil, fmt.Errorf("no default cache store configured")
	}

	store, ok := m.stores[m.defaultStore]
	if !ok {
		return nil, fmt.Errorf("default cache store %s not found", m.defaultStore)
	}
	return store, nil
}

// SetDefault sets the default cache store
func (m *Manager) SetDefault(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.stores[name]; !ok {
		return fmt.Errorf("cache store %s not found", name)
	}

	m.defaultStore = name
	return nil
}
