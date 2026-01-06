package providers

import (
	"fmt"
	"github.com/gin-generator/sugar/config"
	"github.com/gin-generator/sugar/foundation"
	"github.com/gin-generator/sugar/services/storage"
)

// StorageServiceProvider file storage service provider
type StorageServiceProvider struct {
	cfg *config.Config
}

// NewStorageServiceProvider creates a file storage service provider
func NewStorageServiceProvider(cfg *config.Config) *StorageServiceProvider {
	return &StorageServiceProvider{cfg: cfg}
}

// Register registers the service
func (p *StorageServiceProvider) Register(app *foundation.Application) {
	manager := storage.NewManager()
	app.Bind("storage", manager)
}

// Boot boots the service
func (p *StorageServiceProvider) Boot(app *foundation.Application) error {
	service, ok := app.Make("storage")
	if !ok {
		return fmt.Errorf("storage service not found")
	}

	manager := service.(*storage.Manager)

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
