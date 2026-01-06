package providers

import (
	"github.com/gin-generator/sugar/config"
	"github.com/gin-generator/sugar/foundation"
	"github.com/gin-generator/sugar/services/logger"
)

// LoggerServiceProvider logger service provider
type LoggerServiceProvider struct {
	cfg *config.Config
}

// NewLoggerServiceProvider creates a logger service provider
func NewLoggerServiceProvider(cfg *config.Config) *LoggerServiceProvider {
	return &LoggerServiceProvider{cfg: cfg}
}

// Register registers the service
func (p *LoggerServiceProvider) Register(app *foundation.Application) {
	// Logger service is initialized in Boot phase
}

// Boot boots the service
func (p *LoggerServiceProvider) Boot(app *foundation.Application) error {
	var loggerCfg logger.Config
	if err := p.cfg.UnmarshalKey("logger", &loggerCfg); err != nil {
		return err
	}

	log := logger.NewLoggerFromConfig(loggerCfg)
	logger.SetLogger(log)
	app.Bind("logger", log)

	return nil
}

// Name returns the service provider name
func (p *LoggerServiceProvider) Name() string {
	return "Logger"
}
