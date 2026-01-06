package database

import (
	"fmt"
	"gorm.io/gorm"
)

// Global database manager instance
var manager *Manager

// SetManager sets the global database manager
func SetManager(m *Manager) {
	manager = m
}

// DB gets the default database connection (Facade pattern)
func DB() (*gorm.DB, error) {
	if manager == nil {
		return nil, fmt.Errorf("database manager not initialized")
	}
	return manager.DB()
}

// Connection gets a database connection by name (Facade pattern)
func Connection(name string) (*gorm.DB, error) {
	if manager == nil {
		return nil, fmt.Errorf("database manager not initialized")
	}
	return manager.Connection(name)
}
