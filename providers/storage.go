package providers

import (
	"github.com/gin-generator/sugar/foundation"
	"github.com/gin-generator/sugar/services/storage"
)

// StorageServiceProvider file storage service provider
type StorageServiceProvider struct{}

// NewStorageServiceProvider creates a file storage service provider
func NewStorageServiceProvider() *StorageServiceProvider {
	return &StorageServiceProvider{}
}

// Register registers the service
func (p *StorageServiceProvider) Register(app *foundation.Application) {
	manager := storage.NewManager()
	app.Bind(ServiceStorage, manager)
}

// Boot boots the service
func (p *StorageServiceProvider) Boot(app *foundation.Application) error {
	manager := foundation.MustMake[*storage.Manager](app, ServiceStorage)

	// 添加本地存储
	localDisk := storage.NewLocalDisk(storage.LocalConfig{
		Root: "./storage",
	})
	manager.AddDisk("local", localDisk)

	// Add other storage drivers like S3, OSS, etc.

	return nil
}

// Name returns the service provider name
func (p *StorageServiceProvider) Name() string {
	return "Storage"
}
