package storage

import (
	"fmt"
	"io"
	"sync"
)

// Storage file storage interface
type Storage interface {
	Put(path string, contents io.Reader) error
	Get(path string) ([]byte, error)
	Delete(path string) error
	Exists(path string) (bool, error)
	Size(path string) (int64, error)
}

// Manager file storage manager
type Manager struct {
	disks       map[string]Storage
	mu          sync.RWMutex
	defaultDisk string
}

// NewManager creates a new file storage manager
func NewManager() *Manager {
	return &Manager{
		disks: make(map[string]Storage),
	}
}

// AddDisk adds a storage disk
func (m *Manager) AddDisk(name string, disk Storage) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.disks[name] = disk

	if m.defaultDisk == "" {
		m.defaultDisk = name
	}
}

// Disk gets a storage disk by name
func (m *Manager) Disk(name string) (Storage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	disk, ok := m.disks[name]
	if !ok {
		return nil, fmt.Errorf("storage disk %s not found", name)
	}
	return disk, nil
}

// Storage gets the default storage disk
func (m *Manager) Storage() (Storage, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.defaultDisk == "" {
		return nil, fmt.Errorf("no default storage disk configured")
	}

	disk, ok := m.disks[m.defaultDisk]
	if !ok {
		return nil, fmt.Errorf("default storage disk %s not found", m.defaultDisk)
	}
	return disk, nil
}

// SetDefault sets the default disk
func (m *Manager) SetDefault(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.disks[name]; !ok {
		return fmt.Errorf("storage disk %s not found", name)
	}

	m.defaultDisk = name
	return nil
}
