package foundation

import (
	"context"
	"fmt"
	"sync"
)

// Application application container, similar to Laravel's service container
type Application struct {
	context.Context
	mu sync.RWMutex

	// Service container
	services map[string]interface{}

	// Service providers
	providers []ServiceProvider

	// Configuration
	config map[string]interface{}

	// Whether the application has been booted
	booted bool
}

// NewApplication creates a new application instance
func NewApplication() *Application {
	return &Application{
		Context:   context.Background(),
		services:  make(map[string]interface{}),
		providers: make([]ServiceProvider, 0),
		config:    make(map[string]interface{}),
		booted:    false,
	}
}

// Register registers a service provider
func (app *Application) Register(provider ServiceProvider) {
	app.providers = append(app.providers, provider)
	provider.Register(app)
}

// Boot boots all service providers
func (app *Application) Boot() error {
	if app.booted {
		return nil
	}

	for _, provider := range app.providers {
		if err := provider.Boot(app); err != nil {
			return fmt.Errorf("failed to boot provider %s: %w", provider.Name(), err)
		}
	}

	app.booted = true
	return nil
}

// Bind binds a service to the container
func (app *Application) Bind(name string, service interface{}) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.services[name] = service
}

// Make retrieves a service from the container
func (app *Application) Make(name string) (interface{}, bool) {
	app.mu.RLock()
	defer app.mu.RUnlock()
	service, ok := app.services[name]
	return service, ok
}

// MustMake retrieves a service from the container, panics if not found
func (app *Application) MustMake(name string) interface{} {
	service, ok := app.Make(name)
	if !ok {
		panic(fmt.Sprintf("service %s not found in container", name))
	}
	return service
}

// SetConfig sets a configuration value
func (app *Application) SetConfig(key string, value interface{}) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.config[key] = value
}

// GetConfig retrieves a configuration value
func (app *Application) GetConfig(key string) (interface{}, bool) {
	app.mu.RLock()
	defer app.mu.RUnlock()
	value, ok := app.config[key]
	return value, ok
}
