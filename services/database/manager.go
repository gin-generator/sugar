package database

import (
	"fmt"
	"gorm.io/gorm"
	"sync"
)

// Manager database manager for MySQL, PostgreSQL, etc.
type Manager struct {
	connections      map[string]*gorm.DB
	mu               sync.RWMutex
	defaultConnection string
}

// NewManager creates a new database manager
func NewManager() *Manager {
	return &Manager{
		connections: make(map[string]*gorm.DB),
	}
}

// AddConnection adds a database connection
func (m *Manager) AddConnection(name string, db *gorm.DB) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.connections[name] = db

	// First connection becomes the default
	if m.defaultConnection == "" {
		m.defaultConnection = name
	}
}

// Connection gets a database connection by name
func (m *Manager) Connection(name string) (*gorm.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	db, ok := m.connections[name]
	if !ok {
		return nil, fmt.Errorf("database connection %s not found", name)
	}
	return db, nil
}

// DB gets the default database connection
func (m *Manager) DB() (*gorm.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.defaultConnection == "" {
		return nil, fmt.Errorf("no default database connection configured")
	}

	db, ok := m.connections[m.defaultConnection]
	if !ok {
		return nil, fmt.Errorf("default database connection %s not found", m.defaultConnection)
	}
	return db, nil
}

// SetDefault sets the default connection
func (m *Manager) SetDefault(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.connections[name]; !ok {
		return fmt.Errorf("database connection %s not found", name)
	}

	m.defaultConnection = name
	return nil
}
