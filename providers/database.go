package providers

import (
	"fmt"
	"github.com/gin-generator/sugar/foundation"
	"github.com/gin-generator/sugar/services/database"
)

// DatabaseServiceProvider database service provider
type DatabaseServiceProvider struct{}

// NewDatabaseServiceProvider creates a database service provider
func NewDatabaseServiceProvider() *DatabaseServiceProvider {
	return &DatabaseServiceProvider{}
}

// Register registers the service
func (p *DatabaseServiceProvider) Register(app *foundation.Application) {
	manager := database.NewManager()
	app.Bind(ServiceDB, manager)
}

// Boot boots the service
func (p *DatabaseServiceProvider) Boot(app *foundation.Application) error {
	manager := foundation.MustMake[*database.Manager](app, ServiceDB)
	cfg := app.Config

	// Initialize MySQL connections
	for name, mysqlCfg := range cfg.Database.Mysql {
		// Set database name
		mysqlCfg.Database = name

		db, err := database.NewMysqlConnection(name, mysqlCfg)
		if err != nil {
			return fmt.Errorf("failed to connect mysql %s: %w", name, err)
		}

		manager.AddConnection(name, db)
	}

	// Initialize PostgresSQL connections
	for name, pgsqlCfg := range cfg.Database.Pgsql {
		db, err := database.NewPgsqlConnection(name, pgsqlCfg)
		if err != nil {
			return fmt.Errorf("failed to connect pgsql %s: %w", name, err)
		}

		manager.AddConnection(name, db)
	}

	// Set global Facade
	database.SetManager(manager)

	return nil
}

// Name returns the service provider name
func (p *DatabaseServiceProvider) Name() string {
	return "Database"
}
