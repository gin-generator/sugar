package providers

import (
	"github.com/gin-generator/sugar/foundation"
	"github.com/gin-generator/sugar/services/logger"
)

// LoggerServiceProvider logger service provider
type LoggerServiceProvider struct{}

// NewLoggerServiceProvider creates a logger service provider
func NewLoggerServiceProvider() *LoggerServiceProvider {
	return &LoggerServiceProvider{}
}

// Register registers the service
func (p *LoggerServiceProvider) Register(app *foundation.Application) {
	// Logger service is initialized in Boot phase
}

// Boot boots the service
func (p *LoggerServiceProvider) Boot(app *foundation.Application) error {
	cfg := app.Config

	log := logger.NewLoggerFromConfig(cfg.Logger)
	logger.SetLogger(log)
	app.Bind(ServiceLogger, log)

	return nil
}

// Name returns the service provider name
func (p *LoggerServiceProvider) Name() string {
	return "Logger"
}
