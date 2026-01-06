package storage

import (
	"io"
	"os"
	"path/filepath"
)

// LocalConfig local storage configuration
type LocalConfig struct {
	Root string // root directory
}

// LocalDisk local file storage
type LocalDisk struct {
	root string
}

// NewLocalDisk creates a local storage disk
func NewLocalDisk(cfg LocalConfig) *LocalDisk {
	return &LocalDisk{
		root: cfg.Root,
	}
}

// Put saves a file
func (l *LocalDisk) Put(path string, contents io.Reader) error {
	fullPath := filepath.Join(l.root, path)

	// Ensure directory exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, contents)
	return err
}

// Get retrieves file contents
func (l *LocalDisk) Get(path string) ([]byte, error) {
	fullPath := filepath.Join(l.root, path)
	return os.ReadFile(fullPath)
}

// Delete removes a file
func (l *LocalDisk) Delete(path string) error {
	fullPath := filepath.Join(l.root, path)
	return os.Remove(fullPath)
}

// Exists checks if a file exists
func (l *LocalDisk) Exists(path string) (bool, error) {
	fullPath := filepath.Join(l.root, path)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Size gets file size
func (l *LocalDisk) Size(path string) (int64, error) {
	fullPath := filepath.Join(l.root, path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}
