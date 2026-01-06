package foundation

// ServiceProvider service provider interface, similar to Laravel's ServiceProvider
type ServiceProvider interface {
	// Register registers services to the container
	Register(app *Application)

	// Boot boots the service
	Boot(app *Application) error

	// Name returns the service provider name
	Name() string
}
