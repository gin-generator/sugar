package queue

import (
	"context"
	"fmt"
	"sync"
)

// Job queue job interface
type Job interface {
	Handle(ctx context.Context) error
}

// Queue queue interface
type Queue interface {
	Push(job Job) error
	Pop() (Job, error)
	Size() (int, error)
}

// Manager queue manager
type Manager struct {
	connections       map[string]Queue
	mu                sync.RWMutex
	defaultConnection string
}

// NewManager creates a new queue manager
func NewManager() *Manager {
	return &Manager{
		connections: make(map[string]Queue),
	}
}

// AddConnection adds a queue connection
func (m *Manager) AddConnection(name string, queue Queue) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.connections[name] = queue

	if m.defaultConnection == "" {
		m.defaultConnection = name
	}
}

// Connection gets a queue connection by name
func (m *Manager) Connection(name string) (Queue, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	queue, ok := m.connections[name]
	if !ok {
		return nil, fmt.Errorf("queue connection %s not found", name)
	}
	return queue, nil
}

// Queue gets the default queue connection
func (m *Manager) Queue() (Queue, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.defaultConnection == "" {
		return nil, fmt.Errorf("no default queue connection configured")
	}

	queue, ok := m.connections[m.defaultConnection]
	if !ok {
		return nil, fmt.Errorf("default queue connection %s not found", m.defaultConnection)
	}
	return queue, nil
}

// SetDefault sets the default connection
func (m *Manager) SetDefault(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.connections[name]; !ok {
		return fmt.Errorf("queue connection %s not found", name)
	}

	m.defaultConnection = name
	return nil
}
