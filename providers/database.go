package providers

import (
	"fmt"
	"github.com/gin-generator/sugar/config"
	"github.com/gin-generator/sugar/foundation"
	"github.com/gin-generator/sugar/services/database"
)

// DatabaseServiceProvider database service provider
type DatabaseServiceProvider struct {
	cfg *config.Config
}

// NewDatabaseServiceProvider creates a database service provider
func NewDatabaseServiceProvider(cfg *config.Config) *DatabaseServiceProvider {
	return &DatabaseServiceProvider{cfg: cfg}
}

// Register registers the service
func (p *DatabaseServiceProvider) Register(app *foundation.Application) {
	manager := database.NewManager()
	app.Bind("db", manager)
}

// Boot boots the service
func (p *DatabaseServiceProvider) Boot(app *foundation.Application) error {
	service, ok := app.Make("db")
	if !ok {
		return fmt.Errorf("database service not found")
	}

	manager := service.(*database.Manager)

	// Initialize MySQL connections
	mysqlConfigs := p.cfg.GetStringMap("database.mysql")
	for name := range mysqlConfigs {
		var mysqlCfg database.MysqlConfig
		if err := p.cfg.UnmarshalKey(fmt.Sprintf("database.mysql.%s", name), &mysqlCfg); err != nil {
			return fmt.Errorf("failed to unmarshal mysql config %s: %w", name, err)
		}

		// Set database name
		mysqlCfg.Database = name

		db, err := database.NewMysqlConnection(name, mysqlCfg)
		if err != nil {
			return fmt.Errorf("failed to connect mysql %s: %w", name, err)
		}

		// Use database name directly without prefix
		manager.AddConnection(name, db)
	}

	// Initialize PostgreSQL connections
	pgsqlConfigs := p.cfg.GetStringMap("database.pgsql")
	for name := range pgsqlConfigs {
		var pgsqlCfg database.PgsqlConfig
		if err := p.cfg.UnmarshalKey(fmt.Sprintf("database.pgsql.%s", name), &pgsqlCfg); err != nil {
			return fmt.Errorf("failed to unmarshal pgsql config %s: %w", name, err)
		}

		db, err := database.NewPgsqlConnection(name, pgsqlCfg)
		if err != nil {
			return fmt.Errorf("failed to connect pgsql %s: %w", name, err)
		}

		// Use database name directly without prefix
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
