package foundation

import (
	"context"
	"fmt"
	"github.com/gin-generator/sugar/config"
	"sync"
)

// Application application container
type Application struct {
	context.Context
	mu sync.RWMutex

	// Service container
	services map[string]any

	// Service providers
	providers []ServiceProvider

	// Configuration (direct access)
	Config *config.Config

	// Whether the application has been booted
	booted bool
}

// NewApplication creates a new application instance
func NewApplication() *Application {
	return &Application{
		Context:   context.Background(),
		services:  make(map[string]any),
		providers: make([]ServiceProvider, 0),
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
func (app *Application) Bind(name string, service any) {
	app.mu.Lock()
	defer app.mu.Unlock()
	app.services[name] = service
}

// Make retrieves a service from the container with generic type
func Make[T any](app *Application, name string) (T, error) {
	app.mu.RLock()
	defer app.mu.RUnlock()

	var zero T
	service, ok := app.services[name]
	if !ok {
		return zero, fmt.Errorf("service %s not found in container", name)
	}

	typed, ok := service.(T)
	if !ok {
		return zero, fmt.Errorf("service %s type mismatch: expected %T, got %T", name, zero, service)
	}

	return typed, nil
}

// MustMake retrieves a service from the container with generic type, panics if not found
// Only use during application bootstrap phase
func MustMake[T any](app *Application, name string) T {
	service, err := Make[T](app, name)
	if err != nil {
		panic(err)
	}
	return service
}
